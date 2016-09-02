package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const api string = "https://app.scaleft.com/v1/teams/"

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	hostname := d.Get("hostname").(string)
	d.SetId(hostname + "_scaleft")
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {

	key_id := os.Getenv("SCALEFT_KEY_ID")
	key_secret := os.Getenv("SCALEFT_KEY_SECRET")
	key_team := os.Getenv("SCALEFT_TEAM")
	project := os.Getenv("SCALEFT_PROJECT")
	hostname := d.Get("hostname").(string)

	log.Printf("[DEBUG] key_id:%s key_secret:%s key_team:%s project:%s hostname:%s", key_id, key_secret, key_team, project, hostname)

	bearer, err := get_token(key_id, key_secret, key_team)
	if err != nil {
		return fmt.Errorf("Error getting token key_id:%s key_team:%s error:%v", key_id, key_team, err)
	}

	list, err := get_servers(bearer, key_team, project)

	if err != nil {
		return fmt.Errorf("Error getting server list. key_team:%s error:%v", key_team, err)
	}

	ids := get_ids_for_hostname(hostname, list)

	if len(ids) == 0 {
		//	return fmt.Errorf("Error, ScaleFT api returned no servers that matched hostname:%s", hostname)
		//      This should not happen, but if it does, it's ok?
		log.Printf("[WARN] No servers matched for Hostname:%s, Team:%s, Project:%s.  We'll keep going though.", hostname, key_team, project)
		return nil
	}

	for _, id := range ids {
		err := delete_server(bearer, key_team, project, id)
		if err != nil {
			log.Printf("[WARN] Failed to delete server with hostname: %s at ScaleFT ID:%s, error:%s", hostname, id, err)
			//              return fmt.Errorf("Error deleting server at id:%s and key_team:%s project: %s error:%v", id, key_team, project, err)
		}
	}

	return nil
}

type Body struct {
	Key_id     string `json:"key_id"`
	Key_secret string `json:"key_secret"`
}

type Bearer struct {
	Bearer_token string `json:"bearer_token"`
}

type Server struct {
	Id              string                 `json:"id"`
	ProjectName     string                 `json:"project_name"`
	Hostname        string                 `json:"hostname"`
	AltNames        []string               `json:"alt_names"`
	AccessAddress   string                 `json:"access_address"`
	OS              string                 `json:"os"`
	RegisteredAt    time.Time              `json:"registered_at"`
	LastSeen        time.Time              `json:"last_seen"`
	CloudProvider   string                 `json:"cloud_provider"`
	SSHHostKeys     []string               `json:"ssh_host_keys"`
	BrokerHostCerts []string               `json:"broker_host_certs"`
	InstanceDetails map[string]interface{} `json:"instance_details"`
	State           string                 `json:"state"`
}

type Servers struct {
	List []*Server `json:"list"`
}

func get_token(key_id string, key_secret string, key_team string) (string, error) {
	p := &Body{key_id, key_secret}
	jsonStr, err := json.Marshal(p)
	url := api + key_team + "/service_token"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "error", fmt.Errorf("Error getting token key_id:%s key_team:%s status:%s error:%v", key_id, key_team, string(resp.Status), err)
	}

	defer resp.Body.Close()
	b := Bearer{}
	json.NewDecoder(resp.Body).Decode(&b)

	return b.Bearer_token, err
}

func get_logs(bearer_token string, key_team string) string {
	client := &http.Client{}
	url := api + key_team + "/audits"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s
}

func get_servers(bearer_token string, key_team string, project string) (Servers, error) {
	client := &http.Client{}
	url := api + key_team + "/projects/" + project + "/servers"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error listing servers: key_team:%s project: %s status:%s error:%v", key_team, project, string(resp.Status), err)
	}

	s := struct {
		List []*Server `json:"list"`
	}{nil}

	json.NewDecoder(resp.Body).Decode(&s)
	return s, err
}

func delete_server(bearer_token string, key_team string, project string, server_id string) error {
	client := &http.Client{}
	url := api + key_team + "/projects/" + project + "/servers/" + server_id
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error deleting server:%s status:%s error:%v", server_id, string(resp.Status), err)
	}
	return nil
}

func get_ids_for_hostname(hostname string, server_list Servers) []string {
	filtered := make([]string, len(server_list.List))
	for i, l := range server_list.List {
		if hostname == l.Hostname {
			filtered[i] = l.Id
		}
	}
	return filtered
}
