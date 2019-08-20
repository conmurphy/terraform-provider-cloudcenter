
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

func resourcePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourcePlanCreate,
		Read:   resourcePlanRead,
		Update: resourcePlanUpdate,
		Delete: resourcePlanDelete,

		Schema: map[string]*schema.Schema{
			"plan_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_name": &schema.Schema{
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
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"monthly_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"node_hour_increment": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"included_bundle_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"price": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
			"one_time_fee": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"annual_fee": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"storage_rate": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"hourly_rate": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"overage_rate": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"overage_limit": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"restricted_to_app_store_only": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bill_to_vendor": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_rollover": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
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
			"number_of_projects": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePlanCreate(d *schema.ResourceData, m interface{}) error {
	plan_name := d.Get("plan_name").(string)
	tenant_id := d.Get("tenant_id").(string)

	d.SetId(tenant_id + ":" + plan_name)

	client := m.(*cloudcenter.Client)

	newPlan := cloudcenter.Plan{

		Name:            d.Get("plan_name").(string),
		Description:     d.Get("description").(string),
		Type:            d.Get("type").(string),
		ShowOnlyToAdmin: d.Get("show_only_to_admin").(bool),
		Price:           d.Get("price").(float64),
		OnetimeFee:      d.Get("one_time_fee").(float64),
		BillToVendor:    d.Get("bill_to_vendor").(bool),
		TenantId:        d.Get("tenant_id").(string),
	}

	plan, err := client.AddPlan(&newPlan)

	if err != nil {
		return errors.New(err.Error())
	}

	return setPlanResourceData(d, plan)
}

func resourcePlanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR PLAN - TENANT ID INCORRECT")
	}

	plan_id_int, err := strconv.Atoi(d.Get("plan_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR PLAN - PLAN ID INCORRECT")
	}

	plan, err := client.GetPlan(tenant_id_int, plan_id_int)

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR PLAN")
	}

	return setPlanResourceData(d, plan)
}

func resourcePlanUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	newPlan := cloudcenter.Plan{

		Id:              d.Get("plan_id").(string),
		Name:            d.Get("plan_name").(string),
		Description:     d.Get("description").(string),
		Type:            d.Get("type").(string),
		ShowOnlyToAdmin: d.Get("show_only_to_admin").(bool),
		Price:           d.Get("price").(float64),
		OnetimeFee:      d.Get("one_time_fee").(float64),
		BillToVendor:    d.Get("bill_to_vendor").(bool),
		TenantId:        d.Get("tenant_id").(string),
	}

	plan, err := client.UpdatePlan(&newPlan)

	if err != nil {
		return errors.New(err.Error())
	}

	return setPlanResourceData(d, plan)
}

func resourcePlanDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR PLAN - TENANT ID INCORRECT")
	}

	plan_id_int, err := strconv.Atoi(d.Get("plan_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR PLAN - PLAN ID INCORRECT")
	}

	err = client.DeletePlan(tenant_id_int, plan_id_int)

	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId("")
	return nil
}

func setPlanResourceData(d *schema.ResourceData, u *cloudcenter.Plan) error {

	if err := d.Set("plan_id", u.Id); err != nil {
		return errors.New("CANNOT SET ID")
	}
	if err := d.Set("resource", u.Resource); err != nil {
		return errors.New("CANNOT SET RESOURCE")
	}
	if err := d.Set("plan_name", u.Name); err != nil {
		return errors.New("CANNOT SET PLAN NAME")
	}
	if err := d.Set("description", u.Description); err != nil {
		return errors.New("CANNOT SET DESCRIPTION")
	}
	if err := d.Set("tenant_id", u.TenantId); err != nil {
		return errors.New("CANNOT SET TENANT ID")
	}
	if err := d.Set("type", u.Type); err != nil {
		return errors.New("CANNOT SET TYPE")
	}
	if err := d.Set("monthly_limit", u.MonthlyLimit); err != nil {
		return errors.New("CANNOT SET MONTHLY LIMIT")
	}
	if err := d.Set("node_hour_increment", u.NodeHourIncrement); err != nil {
		return errors.New("CANNOT SET NODE HOUR INCREMENT")
	}
	if err := d.Set("included_bundle_id", u.IncludedBundleId); err != nil {
		return errors.New("CANNOT SET INCLUDED BUNDLE ID")
	}
	if err := d.Set("price", u.Price); err != nil {
		return errors.New("CANNOT SET PRICE")
	}
	if err := d.Set("one_time_fee", u.OnetimeFee); err != nil {
		return errors.New("CANNOT SET ONE TIME FEE")
	}
	if err := d.Set("annual_fee", u.AnnualFee); err != nil {
		return errors.New("CANNOT SET ANNUAL FEE")
	}
	if err := d.Set("storage_rate", u.StorageRate); err != nil {
		return errors.New("CANNOT SET STORAGE RATE")
	}
	if err := d.Set("hourly_rate", u.HourlyRate); err != nil {
		return errors.New("CANNOT SET HOURLY RATE")
	}
	if err := d.Set("overage_rate", u.OverageRate); err != nil {
		return errors.New("CANNOT SET OVERAGE RATE")
	}
	if err := d.Set("overage_limit", u.OverageLimit); err != nil {
		return errors.New("CANNOT SET OVERAGE LIMIT")
	}
	if err := d.Set("restricted_to_app_store_only", u.RestrictedToAppStoreOnly); err != nil {
		return errors.New("CANNOT SET VALUE - RESTRICTED TO APP STORE ONLY")
	}
	if err := d.Set("bill_to_vendor", u.BillToVendor); err != nil {
		return errors.New("CANNOT SET BILL TO VENDOR")
	}
	if err := d.Set("enable_rollover", u.EnableRollover); err != nil {
		return errors.New("CANNOT SET VALUE - ENABLE ROLLOVER")
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
	if err := d.Set("number_of_projects", u.NumberOfProjects); err != nil {
		return errors.New("CANNOT SET NUMBER OF PROJECTS")
	}
	return nil
}
