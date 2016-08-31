package main

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"net/http"
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

	key_id := d.Get("key_id").(string)
	key_secret := d.Get("key_secret").(string)
	key_team := d.Get("team").(string)
	project := d.Get("project").(string)
	hostname := d.Get("hostname").(string)

	bearer := get_token(key_id, key_secret, key_team)
	list := get_servers(bearer, key_team, project)

	ids := get_ids_for_hostname(hostname, list)

	for _, id := range ids {
		delete_server(bearer, key_team, project, id)
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

func get_token(key_id string, key_secret string, key_team string) string {
	p := &Body{key_id, key_secret}
	jsonStr, err := json.Marshal(p)
	url := api + key_team + "/service_token"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	b := Bearer{}
	json.NewDecoder(resp.Body).Decode(&b)

	return b.Bearer_token
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

func get_servers(bearer_token string, key_team string, project string) Servers {
	client := &http.Client{}
	url := api + key_team + "/projects/" + project + "/servers"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	s := struct {
		List []*Server `json:"list"`
	}{nil}

	json.NewDecoder(resp.Body).Decode(&s)
	return s
}

func delete_server(bearer_token string, key_team string, project string, server_id string) string {
	client := &http.Client{}
	url := api + key_team + "/projects/" + project + "/servers/" + server_id
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	s := string(resp.Status)
	return s
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
