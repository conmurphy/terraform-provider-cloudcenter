package main

import (
	"errors"
	"github.com/cloudcenter-clientlibrary-go/cloudcenter"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func resourceBundle() *schema.Resource {
	return &schema.Resource{
		Create: resourceBundleCreate,
		Read:   resourceBundleRead,
		Update: resourceBundleUpdate,
		Delete: resourceBundleDelete,

		Schema: map[string]*schema.Schema{
			"bundle_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bundle_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"limit": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
			"price": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
			"expiration_date": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
			"expiration_months": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_only_to_admin": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"number_of_users": &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceBundleCreate(d *schema.ResourceData, m interface{}) error {

	bundle_name := d.Get("bundle_name").(string)
	tenant_id := d.Get("tenant_id").(string)

	d.SetId(tenant_id + ":" + bundle_name)

	client := m.(*cloudcenter.Client)

	newBundle := cloudcenter.Bundle{

		Name:           d.Get("bundle_name").(string),
		Type:           d.Get("type").(string),
		Limit:          d.Get("limit").(float64),
		Price:          d.Get("price").(float64),
		ExpirationDate: d.Get("expiration_date").(float64),
		TenantId:       d.Get("tenant_id").(string),
	}

	bundle, err := client.AddBundle(&newBundle)

	if err != nil {
		return errors.New(err.Error())
	}

	return setBundleResourceData(d, bundle)
}

func resourceBundleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR BUNDLE - TENANT ID INCORRECT")
	}

	bundle_id_int, err := strconv.Atoi(d.Get("bundle_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR BUNDLE - BUNDLE ID INCORRECT")
	}

	bundle, err := client.GetBundle(tenant_id_int, bundle_id_int)

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR BUNDLE")
	}

	return setBundleResourceData(d, bundle)
}

func resourceBundleUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	newBundle := cloudcenter.Bundle{

		Id:             d.Get("bundle_id").(string),
		Name:           d.Get("bundle_name").(string),
		Type:           d.Get("type").(string),
		Limit:          d.Get("limit").(float64),
		Price:          d.Get("price").(float64),
		ExpirationDate: d.Get("expiration_date").(float64),
		TenantId:       d.Get("tenant_id").(string),
	}

	bundle, err := client.UpdateBundle(&newBundle)

	if err != nil {
		return errors.New(err.Error())
	}

	return setBundleResourceData(d, bundle)
}

func resourceBundleDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	tenant_id_int, err := strconv.Atoi(d.Get("tenant_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR BUNDLE - TENANT ID INCORRECT")
	}

	bundle_id_int, err := strconv.Atoi(d.Get("bundle_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR BUNDLE - BUNDLE ID INCORRECT")
	}

	err = client.DeleteBundle(tenant_id_int, bundle_id_int)

	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId("")
	return nil
}

func setBundleResourceData(d *schema.ResourceData, u *cloudcenter.Bundle) error {

	if err := d.Set("bundle_id", u.Id); err != nil {
		return errors.New("CANNOT SET ID")
	}
	if err := d.Set("bundle_name", u.Name); err != nil {
		return errors.New("CANNOT SET BUNDLE NAME")
	}
	if err := d.Set("resource", u.Resource); err != nil {
		return errors.New("CANNOT SET RESOURCE")
	}
	if err := d.Set("type", u.Type); err != nil {
		return errors.New("CANNOT SET TYPE")
	}
	if err := d.Set("limit", u.Limit); err != nil {
		return errors.New("CANNOT SET LIMIT")
	}
	if err := d.Set("price", u.Price); err != nil {
		return errors.New("CANNOT SET PRICE")
	}
	if err := d.Set("expiration_date", u.ExpirationDate); err != nil {
		return errors.New("CANNOT SET EXPIRATION DATE")
	}
	if err := d.Set("expiration_months", u.ExpirationMonths); err != nil {
		return errors.New("CANNOT SET EXPIRATION MONTHS")
	}
	if err := d.Set("disabled", u.Disabled); err != nil {
		return errors.New("CANNOT SET DISABLED")
	}
	if err := d.Set("show_only_to_admin", u.ShowOnlyToAdmin); err != nil {
		return errors.New("CANNOT SET SHOW ONLY TO ADMIN FIELD")
	}
	if err := d.Set("number_of_users", u.NumberOfUsers); err != nil {
		return errors.New("CANNOT SET NUMBER OF USERS")
	}
	if err := d.Set("tenant_id", u.TenantId); err != nil {
		return errors.New("CANNOT SET TENANT ID")
	}

	return nil
}
