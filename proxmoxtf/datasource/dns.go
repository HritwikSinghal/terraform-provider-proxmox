/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentDNSDomain   = "domain"
	mkDataSourceVirtualEnvironmentDNSNodeName = "node_name"
	mkDataSourceVirtualEnvironmentDNSServers  = "servers"
)

// DNS returns a resource for DNS settings on a node.
func DNS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentDNSDomain: {
				Type:        schema.TypeString,
				Description: "The DNS search domain",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentDNSNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentDNSServers: {
				Type:        schema.TypeList,
				Description: "The DNS servers",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: dnsRead,
	}
}

func dnsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentDNSNodeName).(string)

	dns, err := api.Node(nodeName).GetDNS(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s_dns", nodeName))

	if dns.SearchDomain != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentDNSDomain, *dns.SearchDomain)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentDNSDomain, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	var servers []interface{}

	if dns.Server1 != nil {
		servers = append(servers, *dns.Server1)
	}

	if dns.Server2 != nil {
		servers = append(servers, *dns.Server2)
	}

	if dns.Server3 != nil {
		servers = append(servers, *dns.Server3)
	}

	err = d.Set(mkDataSourceVirtualEnvironmentDNSServers, servers)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
