# terraform-provider-scaleft

For now, just support deleting all instances of a server that exist for the hostname of a server that is about to be destroyed.  The goal is to cleanup some of our very short lived services.

No one in there right mind would use this.  I have very little clue what I am doing.

The create does nothing other than record and object in the tfstate so that we know on destroy we need to delete something.  It would probably be better to grab the ID post enrollment, but then we'd had to rely on ScaleFTs api being available to deploy.  Might be good to try to get the ID but be ok with the api failing and still fall back to what it does now, search for all IDs for a hostname and remove each of them.  For our use case, this is always a sane option, but it may not be for everyone elses.

Example use:
```
provider "scaleft" {}

resource "scaleft_server" "machine" {
  hostname = "${var.tier}${count.index}${var.aws_suffix}.${var.env_name}.${var.internal_tld}",
  count = "${var.aws_count}"
}
```
