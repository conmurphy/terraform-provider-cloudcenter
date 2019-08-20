
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

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,

		Schema: map[string]*schema.Schema{
			"role_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"perms": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"oob_role": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"object_permissions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"permissions": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"groups": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": &schema.Schema{
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

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {

	role_name := d.Get("role_name").(string)
	tenant_id := d.Get("tenant_id").(string)

	d.SetId(tenant_id + ":" + role_name)

	client := m.(*cloudcenter.Client)

	var objectPerms []cloudcenter.ObjectPerm

	allObjectPerms := d.Get("object_permissions").([]interface{})

	for _, object := range allObjectPerms {

		o, _ := object.(map[string]interface{})

		perms := []string{}

		permissions := o["permissions"].([]interface{})

		for _, permission := range permissions {
			perms = append(perms, permission.(string))
		}

		newObjectPerm := cloudcenter.ObjectPerm{
			ObjectType: o["object_type"].(string),
			Perms:      perms,
		}

		objectPerms = append(objectPerms, newObjectPerm)

	}

	var users []cloudcenter.User

	allUsers := d.Get("users").([]interface{})

	for _, user := range allUsers {

		u, _ := user.(map[string]interface{})

		newUser := cloudcenter.User{
			Id: u["user_id"].(string),
		}

		users = append(users, newUser)

	}

	var groups []cloudcenter.Group

	allGroups := d.Get("groups").([]interface{})

	for _, group := range allGroups {

		g, _ := group.(map[string]interface{})

		newGroup := cloudcenter.Group{
			Id: g["group_id"].(string),
		}

		groups = append(groups, newGroup)

	}

	perms := []string{}
	for _, perm := range d.Get("perms").([]interface{}) {
		perms = append(perms, perm.(string))
	}

	newRole := cloudcenter.Role{

		TenantId:    d.Get("tenant_id").(string),
		Name:        d.Get("role_name").(string),
		Description: d.Get("description").(string),
		ObjectPerms: objectPerms,
		Users:       users,
		Groups:      groups,
		Perms:       perms,
	}

	role, err := client.AddRole(&newRole)

	if err != nil {
		return errors.New(err.Error())
	}

	return setRoleResourceData(d, role)
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR ROLE - TENANT ID INCORRECT")
	}

	role_id_int, err := strconv.Atoi(d.Get("role_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR ROLE - ROLE ID INCORRECT")
	}

	role, err := client.GetRole(tenant_id_int, role_id_int)

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR ROLE")
	}

	return setRoleResourceData(d, role)
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	var objectPerms []cloudcenter.ObjectPerm

	allObjectPerms := d.Get("object_permissions").([]interface{})

	for _, object := range allObjectPerms {

		o, _ := object.(map[string]interface{})

		perms := []string{}

		permissions := o["permissions"].([]interface{})

		for _, permission := range permissions {
			perms = append(perms, permission.(string))
		}

		newObjectPerm := cloudcenter.ObjectPerm{
			ObjectType: o["object_type"].(string),
			Perms:      perms,
		}

		objectPerms = append(objectPerms, newObjectPerm)

	}

	var users []cloudcenter.User

	allUsers := d.Get("users").([]interface{})

	for _, user := range allUsers {

		u, _ := user.(map[string]interface{})

		newUser := cloudcenter.User{
			Id: u["user_id"].(string),
		}

		users = append(users, newUser)

	}

	var groups []cloudcenter.Group

	allGroups := d.Get("groups").([]interface{})

	for _, group := range allGroups {

		g, _ := group.(map[string]interface{})

		newGroup := cloudcenter.Group{
			Id: g["group_id"].(string),
		}

		groups = append(groups, newGroup)

	}

	perms := []string{}
	for _, perm := range d.Get("perms").([]interface{}) {
		perms = append(perms, perm.(string))
	}

	newRole := cloudcenter.Role{

		Id:          d.Get("role_id").(string),
		TenantId:    d.Get("tenant_id").(string),
		Name:        d.Get("role_name").(string),
		Description: d.Get("description").(string),
		ObjectPerms: objectPerms,
		Users:       users,
		Groups:      groups,
		Perms:       perms,
	}

	role, err := client.UpdateRole(&newRole)

	if err != nil {
		return errors.New(err.Error())
	}

	return setRoleResourceData(d, role)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR ROLE - TENANT ID INCORRECT")
	}

	role_id_int, err := strconv.Atoi(d.Get("role_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR ROLE - ROLE ID INCORRECT")
	}

	err = client.DeleteRole(tenant_id_int, role_id_int)

	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId("")
	return nil
}

func setRoleResourceData(d *schema.ResourceData, u *cloudcenter.Role) error {

	if err := d.Set("role_id", u.Id); err != nil {
		return errors.New("CANNOT SET ID")
	}
	if err := d.Set("role_name", u.Name); err != nil {
		return errors.New("CANNOT SET ROLE NAME")
	}
	if err := d.Set("perms", u.Perms); err != nil {
		return errors.New("CANNOT SET PERMS")
	}
	if err := d.Set("resource", u.Resource); err != nil {
		return errors.New("CANNOT SET RESOURCE")
	}
	if err := d.Set("description", u.Description); err != nil {
		return errors.New("CANNOT SET DESCRIPTION")
	}
	if err := d.Set("object_permissions", d.Get("object_permissions").([]interface{})); err != nil {
		return errors.New("CANNOT SET OBJECT PERMISSIONS")
	}
	if err := d.Set("users", d.Get("users").([]interface{})); err != nil {
		return errors.New("CANNOT SET USERS")
	}
	if err := d.Set("groups", d.Get("groups").([]interface{})); err != nil {
		return errors.New("CANNOT SET GROUPS")
	}
	if err := d.Set("created", u.Created); err != nil {
		return errors.New("CANNOT SET CREATED VALUE")
	}
	if err := d.Set("last_updated", u.LastUpdated); err != nil {
		return errors.New("CANNOT SET LAST UPDATED VALUE")
	}
	if err := d.Set("oob_role", u.OobRole); err != nil {
		return errors.New("CANNOT SET CREATED BY SSO VALUE")
	}
	if err := d.Set("tenant_id", u.TenantId); err != nil {
		return errors.New("CANNOT SET TENANT ID")
	}

	return nil
}
