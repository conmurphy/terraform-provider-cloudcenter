package main

import (
	"errors"
	"github.com/cloudcenter-clientlibrary-go/cloudcenter"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func resourceActivationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceActivationProfileCreate,
		Read:   resourceActivationProfileRead,
		Update: resourceActivationProfileUpdate,
		Delete: resourceActivationProfileDelete,

		Schema: map[string]*schema.Schema{
			"activation_profile_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"activation_profile_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"bundle_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"contract_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_environment_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"activate_regions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"agree_to_contract": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"send_activation_email": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceActivationProfileCreate(d *schema.ResourceData, m interface{}) error {

	activation_profile_name := d.Get("activation_profile_name").(string)
	tenant_id := strconv.Itoa(d.Get("tenant_id").(int))

	d.SetId(tenant_id + ":" + activation_profile_name)

	client := m.(*cloudcenter.Client)

	var activateRegions []cloudcenter.ActivateRegion

	allRegions := d.Get("activate_regions").([]interface{})

	for _, region := range allRegions {

		r, _ := region.(map[string]interface{})

		newActivateRegion := cloudcenter.ActivateRegion{
			RegionId: r["region_id"].(string),
		}

		activateRegions = append(activateRegions, newActivateRegion)

	}

	newActivationProfile := cloudcenter.ActivationProfile{

		Name:        d.Get("activation_profile_name").(string),
		Description: d.Get("description").(string),
		TenantId:    d.Get("tenant_id").(int),
		PlanId:      d.Get("plan_id").(string),
		BundleId:    d.Get("bundle_id").(string),
		ContractId:  d.Get("contract_id").(string),
		DepEnvId:    d.Get("deployment_environment_id").(string),
		//ImportApps:          d.Get("import_apps").([]string),
		ActivateRegions:     activateRegions,
		AgreeToContract:     d.Get("agree_to_contract").(bool),
		SendActivationEmail: d.Get("send_activation_email").(bool),
	}

	activationProfile, err := client.AddActivationProfile(&newActivationProfile)

	if err != nil {
		return errors.New(err.Error())
	}

	return setActivationProfileResourceData(d, activationProfile)
}

func resourceActivationProfileRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudcenter.Client)

	tenant_id_int := d.Get("tenant_id").(int)

	activation_profile_id_int, err := strconv.Atoi(d.Get("activation_profile_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR ACTIVATION PROFILE - ACTIVATION PROFILE ID INCORRECT OR NOT FOUND")
	}

	activationProfile, err := client.GetActivationProfile(tenant_id_int, activation_profile_id_int)

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR ACTIVATION PROFILE")
	}

	return setActivationProfileResourceData(d, activationProfile)
}

func resourceActivationProfileUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	var activateRegions []cloudcenter.ActivateRegion

	allRegions := d.Get("activate_regions").([]interface{})

	for _, region := range allRegions {

		r, _ := region.(map[string]interface{})

		newActivateRegion := cloudcenter.ActivateRegion{
			RegionId: r["region_id"].(string),
		}

		activateRegions = append(activateRegions, newActivateRegion)

	}

	newActivationProfile := cloudcenter.ActivationProfile{

		Id:          d.Get("activation_profile_id").(string),
		Name:        d.Get("activation_profile_name").(string),
		Description: d.Get("description").(string),
		TenantId:    d.Get("tenant_id").(int),
		PlanId:      d.Get("plan_id").(string),
		BundleId:    d.Get("bundle_id").(string),
		ContractId:  d.Get("contract_id").(string),
		DepEnvId:    d.Get("deployment_environment_id").(string),
		//ImportApps:          d.Get("import_apps").([]string),
		ActivateRegions:     activateRegions,
		AgreeToContract:     d.Get("agree_to_contract").(bool),
		SendActivationEmail: d.Get("send_activation_email").(bool),
	}

	activationProfile, err := client.UpdateActivationProfile(&newActivationProfile)

	if err != nil {
		return errors.New(err.Error())
	}

	return setActivationProfileResourceData(d, activationProfile)
}

func resourceActivationProfileDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	tenant_id_int := d.Get("tenant_id").(int)

	activation_profile_id_int, err := strconv.Atoi(d.Get("activation_profile_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR ACTIVATION PROFILE - ACTIVATION PROFILE ID INCORRECT OR NOT FOUND")
	}

	err = client.DeleteActivationProfile(tenant_id_int, activation_profile_id_int)

	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId("")
	return nil
}

func setActivationProfileResourceData(d *schema.ResourceData, u *cloudcenter.ActivationProfile) error {

	if err := d.Set("activation_profile_id", u.Id); err != nil {
		return errors.New("CANNOT SET ID")
	}
	if err := d.Set("activation_profile_name", u.Name); err != nil {
		return errors.New("CANNOT SET ACTIVATION PROFILE NAME")
	}
	if err := d.Set("description", u.Description); err != nil {
		return errors.New("CANNOT SET DESCRIPTION")
	}
	if err := d.Set("tenant_id", u.TenantId); err != nil {
		return errors.New("CANNOT SET TENANT ID")
	}
	if err := d.Set("plan_id", u.PlanId); err != nil {
		return errors.New("CANNOT SET PLAN ID")
	}
	if err := d.Set("bundle_id", u.BundleId); err != nil {
		return errors.New("CANNOT SET BUNDLE ID")
	}
	if err := d.Set("contract_id", u.ContractId); err != nil {
		return errors.New("CANNOT SET CONTRACT ID")
	}
	if err := d.Set("deployment_environment_id", u.DepEnvId); err != nil {
		return errors.New("CANNOT SET DEPLOYMENT ENVIRONMENT ID")
	}

	if err := d.Set("activate_regions", d.Get("activate_regions").([]interface{})); err != nil {
		return errors.New("CANNOT SET VALUE - ACTIVATE REGIONS")
	}
	if err := d.Set("agree_to_contract", u.AgreeToContract); err != nil {
		return errors.New("CANNOT SET VALUE - AGREE TO CONTRACT")
	}
	if err := d.Set("send_activation_email", u.SendActivationEmail); err != nil {
		return errors.New("CANNOT SET VALUE - SEND ACTIVATION EMAIL")
	}
	//if err := d.Set("import_apps", u.ImportApps); err != nil {
	//	return errors.New("CANNOT SET VALUE - IMPORT APPS")
	//}

	return nil
}
