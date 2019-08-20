
/*
Copyright (c) 2019 Cisco and/or its affiliates.

This software is licensed to you under the terms of the Cisco Sample
Code License, Version 1.0 (the "License"). You may obtain a copy of the
License at

               https://developer.cisco.com/docs/licenses

All use of the material herein must be in accordance with the terms of
the License. All rights not expressly granted by the License are
reserved. Unless required by applicable law or agreed to separately in
writing, software distributed under the License is distributed on an "AS
IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied.

*/

package main

import (
	"errors"
	"github.com/cloudcenter-clientlibrary-go/cloudcenter"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"created": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_by_sso": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"roles": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {

	group_name := d.Get("group_name").(string)
	tenant_id := d.Get("tenant_id").(string)

	d.SetId(tenant_id + ":" + group_name)

	client := m.(*cloudcenter.Client)

	var users []cloudcenter.User

	allUsers := d.Get("users").([]interface{})

	for _, user := range allUsers {

		u, _ := user.(map[string]interface{})

		newUser := cloudcenter.User{
			Id: u["user_id"].(string),
		}

		users = append(users, newUser)

	}

	var roles []cloudcenter.Role

	allRoles := d.Get("roles").([]interface{})

	for _, role := range allRoles {

		r, _ := role.(map[string]interface{})

		newRole := cloudcenter.Role{
			Id: r["role_id"].(string),
		}

		roles = append(roles, newRole)

	}

	newGroup := cloudcenter.Group{

		TenantId:    d.Get("tenant_id").(string),
		Name:        d.Get("group_name").(string),
		Description: d.Get("description").(string),
		Users:       users,
		Roles:       roles,
	}

	group, err := client.AddGroup(&newGroup)

	if err != nil {
		return errors.New(err.Error())
	}

	return setGroupResourceData(d, group)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR GROUP - TENANT ID INCORRECT")
	}

	group_id_int, err := strconv.Atoi(d.Get("group_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR GROUP - GROUP ID INCORRECT")
	}

	group, err := client.GetGroup(tenant_id_int, group_id_int)

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR GROUP")
	}

	return setGroupResourceData(d, group)
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	var users []cloudcenter.User

	allUsers := d.Get("users").([]interface{})

	for _, user := range allUsers {

		u, _ := user.(map[string]interface{})

		newUser := cloudcenter.User{
			Id: u["user_id"].(string),
		}

		users = append(users, newUser)

	}

	var roles []cloudcenter.Role

	allRoles := d.Get("roles").([]interface{})

	for _, role := range allRoles {

		r, _ := role.(map[string]interface{})

		newRole := cloudcenter.Role{
			Id: r["role_id"].(string),
		}

		roles = append(roles, newRole)

	}

	newGroup := cloudcenter.Group{

		Id:          d.Get("group_id").(string),
		TenantId:    d.Get("tenant_id").(string),
		Name:        d.Get("group_name").(string),
		Description: d.Get("description").(string),
		Users:       users,
		Roles:       roles,
	}

	group, err := client.UpdateGroup(&newGroup)

	if err != nil {
		return errors.New(err.Error())
	}

	return setGroupResourceData(d, group)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR GROUP - TENANT ID INCORRECT")
	}

	group_id_int, err := strconv.Atoi(d.Get("group_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR GROUP - GROUP ID INCORRECT")
	}

	err = client.DeleteGroup(tenant_id_int, group_id_int)

	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId("")
	return nil
}

func setGroupResourceData(d *schema.ResourceData, u *cloudcenter.Group) error {

	if err := d.Set("group_id", u.Id); err != nil {
		return errors.New("CANNOT SET ID")
	}
	if err := d.Set("group_name", u.Name); err != nil {
		return errors.New("CANNOT SET GROUP NAME")
	}
	if err := d.Set("resource", u.Resource); err != nil {
		return errors.New("CANNOT SET RESOURCE")
	}
	if err := d.Set("description", u.Description); err != nil {
		return errors.New("CANNOT SET DESCRIPTION")
	}
	if err := d.Set("users", d.Get("users").([]interface{})); err != nil {
		return errors.New("CANNOT SET USERS")
	}
	if err := d.Set("roles", d.Get("roles").([]interface{})); err != nil {
		return errors.New("CANNOT SET ROLES")
	}
	if err := d.Set("created", u.Created); err != nil {
		return errors.New("CANNOT SET CREATED VALUE")
	}
	if err := d.Set("last_updated", u.LastUpdated); err != nil {
		return errors.New("CANNOT SET LAST UPDATED VALUE")
	}
	if err := d.Set("created_by_sso", u.CreatedBySso); err != nil {
		return errors.New("CANNOT SET CREATED BY SSO VALUE")
	}
	if err := d.Set("tenant_id", u.TenantId); err != nil {
		return errors.New("CANNOT SET TENANT ID")
	}

	return nil
}
