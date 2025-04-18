package github

import (
	"context"

	"github.com/google/go-github/v66/github"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGithubCodespacesOrganizationSecrets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGithubCodespacesOrganizationSecretsRead,

		Schema: map[string]*schema.Schema{
			"secrets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"visibility": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGithubCodespacesOrganizationSecretsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	owner := meta.(*Owner).name

	options := github.ListOptions{
		PerPage: 100,
	}

	var all_secrets []map[string]string
	for {
		secrets, resp, err := client.Codespaces.ListOrgSecrets(context.TODO(), owner, &options)
		if err != nil {
			return err
		}
		for _, secret := range secrets.Secrets {
			new_secret := map[string]string{
				"name":       secret.Name,
				"visibility": secret.Visibility,
				"created_at": secret.CreatedAt.String(),
				"updated_at": secret.UpdatedAt.String(),
			}
			all_secrets = append(all_secrets, new_secret)

		}
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	d.SetId(owner)
	err := d.Set("secrets", all_secrets)
	if err != nil {
		return err
	}

	return nil
}
