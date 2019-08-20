
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

func resourceContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceContractCreate,
		Read:   resourceContractRead,
		Update: resourceContractUpdate,
		Delete: resourceContractDelete,

		Schema: map[string]*schema.Schema{
			"contract_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"contract_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"length": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"terms": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"discount_rate": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"show_only_to_admin": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"number_of_users": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceContractCreate(d *schema.ResourceData, m interface{}) error {
	contract_name := d.Get("contract_name").(string)
	tenant_id := d.Get("tenant_id").(string)

	d.SetId(tenant_id + ":" + contract_name)

	client := m.(*cloudcenter.Client)

	newContract := cloudcenter.Contract{

		Name:         d.Get("contract_name").(string),
		Description:  d.Get("description").(string),
		Length:       d.Get("length").(int),
		Terms:        d.Get("terms").(string),
		DiscountRate: d.Get("discount_rate").(float64),
		TenantId:     d.Get("tenant_id").(string),
	}

	contract, err := client.AddContract(&newContract)

	if err != nil {
		return errors.New(err.Error())
	}

	return setContractResourceData(d, contract)
}

func resourceContractRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR CONTRACT - TENANT ID INCORRECT")
	}

	contract_id_int, err := strconv.Atoi(d.Get("contract_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR CONTRACT - CONTRACT ID INCORRECT")
	}

	contract, err := client.GetContract(tenant_id_int, contract_id_int)

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR CONTRACT")
	}

	return setContractResourceData(d, contract)
}

func resourceContractUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	newContract := cloudcenter.Contract{

		Id:           d.Get("contract_id").(string),
		Name:         d.Get("contract_name").(string),
		Description:  d.Get("description").(string),
		Length:       d.Get("length").(int),
		Terms:        d.Get("terms").(string),
		DiscountRate: d.Get("discount_rate").(float64),
		TenantId:     d.Get("tenant_id").(string),
	}

	contract, err := client.UpdateContract(&newContract)

	if err != nil {
		return errors.New(err.Error())
	}

	return setContractResourceData(d, contract)
}

func resourceContractDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR CONTRACT - TENANT ID INCORRECT")
	}

	contact_id_int, err := strconv.Atoi(d.Get("contract_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR CONTRACT - CONTRACT ID INCORRECT")
	}

	err = client.DeleteContract(tenant_id_int, contact_id_int)

	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId("")
	return nil
}

func setContractResourceData(d *schema.ResourceData, u *cloudcenter.Contract) error {

	if err := d.Set("contract_id", u.Id); err != nil {
		return errors.New("CANNOT SET ID")
	}
	if err := d.Set("resource", u.Resource); err != nil {
		return errors.New("CANNOT SET RESOURCE")
	}
	if err := d.Set("contract_name", u.Name); err != nil {
		return errors.New("CANNOT SET CONTRACT NAME")
	}
	if err := d.Set("description", u.Description); err != nil {
		return errors.New("CANNOT SET DESCRIPTION")
	}
	if err := d.Set("tenant_id", u.TenantId); err != nil {
		return errors.New("CANNOT SET TENANT ID")
	}
	if err := d.Set("length", u.Length); err != nil {
		return errors.New("CANNOT SET LENGTH")
	}
	if err := d.Set("terms", u.Terms); err != nil {
		return errors.New("CANNOT SET TERMS")
	}
	if err := d.Set("discount_rate", u.DiscountRate); err != nil {
		return errors.New("CANNOT SET DISCOUNT RATE")
	}
	if err := d.Set("disabled", u.Disabled); err != nil {
		return errors.New("CANNOT SET DISABLED")
	}
	if err := d.Set("show_only_to_admin", u.ShowOnlyToAdmin); err != nil {
		return errors.New("CANNOT SET VALUE - SHOW ONLY TO ADMIN")
	}
	if err := d.Set("number_of_users", u.NumberOfUsers); err != nil {
		return errors.New("CANNOT SET NUMBER OF USERS")
	}
	return nil
}
