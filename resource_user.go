package main

import (
	"errors"
	"github.com/cloudcenter-clientlibrary-go/cloudcenter"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"email_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"company_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_verified": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"phone_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_keys": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"disable_reason": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_source": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"detail": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"activation_data": &schema.Schema{
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
			"co_admin": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tenant_admin": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"activation_profile_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_subscription_plan": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	email_address := d.Get("email_address").(string)
	d.SetId(email_address)

	client := m.(*cloudcenter.Client)

	newUser := cloudcenter.User{

		FirstName:   d.Get("first_name").(string),
		LastName:    d.Get("last_name").(string),
		Password:    d.Get("password").(string),
		EmailAddr:   d.Get("email_address").(string),
		CompanyName: d.Get("company_name").(string),
		PhoneNumber: d.Get("phone_number").(string),
		TenantId:    d.Get("tenant_id").(string),
	}

	user, err := client.AddUser(&newUser)

	if err != nil {
		return errors.New(err.Error())
	}

	return setUserResourceData(d, user)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudcenter.Client)
	user, err := client.GetUserFromEmail(d.Get("email_address").(string))
	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR USER: " + d.Get("username").(string))
	}

	return setUserResourceData(d, user)
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	newUser := cloudcenter.User{
		Id:            d.Get("user_id").(string),
		FirstName:     d.Get("first_name").(string),
		LastName:      d.Get("last_name").(string),
		Password:      d.Get("password").(string),
		EmailAddr:     d.Get("email_address").(string),
		CompanyName:   d.Get("company_name").(string),
		PhoneNumber:   d.Get("phone_number").(string),
		TenantId:      d.Get("tenant_id").(string),
		Username:      d.Get("username").(string),
		AccountSource: d.Get("account_source").(string),
		Type:          d.Get("type").(string),
	}

	user, err := client.UpdateUser(&newUser)

	if err != nil {
		return errors.New(err.Error())
	}

	return setUserResourceData(d, user)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	err := client.DeleteUserByEmail(d.Get("email_address").(string))

	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId("")
	return nil
}

func setUserResourceData(d *schema.ResourceData, u *cloudcenter.User) error {

	if err := d.Set("user_id", u.Id); err != nil {
		return errors.New("CANNOT SET ID")
	}
	if err := d.Set("username", u.Username); err != nil {
		return errors.New("CANNOT SET USERNAME")
	}
	if err := d.Set("email_address", u.EmailAddr); err != nil {
		return errors.New("CANNOT SET EMAIL")
	}
	if err := d.Set("first_name", u.FirstName); err != nil {
		return errors.New("CANNOT SET FIRST NAME")
	}
	if err := d.Set("last_name", u.LastName); err != nil {
		return errors.New("CANNOT SET LAST NAME")
	}
	if err := d.Set("tenant_id", u.TenantId); err != nil {
		return errors.New("CANNOT SET TENANT ID")
	}
	if err := d.Set("type", u.Type); err != nil {
		return errors.New("CANNOT SET TYPE")
	}
	if err := d.Set("account_source", u.AccountSource); err != nil {
		return errors.New("CANNOT SET ACCOUNT SOURCE")
	}
	return nil
}
