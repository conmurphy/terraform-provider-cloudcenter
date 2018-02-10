package main

import (
	"errors"
	"github.com/cloudcenter-clientlibrary-go/cloudcenter"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageCreate,
		Read:   resourceImageRead,
		Update: resourceImageUpdate,
		Delete: resourceImageDelete,

		Schema: map[string]*schema.Schema{
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"internal_image_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"visibility": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			/*"tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},*/
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_image": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"num_of_nics": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"attach_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceImageCreate(d *schema.ResourceData, m interface{}) error {

	image_name := d.Get("image_name").(string)
	tenant_id := d.Get("tenant_id").(int)

	d.SetId(strconv.Itoa(tenant_id) + ":" + image_name)

	client := m.(*cloudcenter.Client)

	newImage := cloudcenter.Image{

		Name:        d.Get("image_name").(string),
		Description: d.Get("description").(string),
		Visibility:  d.Get("visibility").(string),
		NumOfNICs:   d.Get("num_of_nics").(int),
		OSName:      d.Get("os_name").(string),
		Enabled:     d.Get("enabled").(bool),
		TenantId:    d.Get("tenant_id").(int),
		ImageType:   d.Get("image_type").(string),
	}

	image, err := client.AddImage(&newImage)

	if err != nil {
		return errors.New(err.Error())
	}

	return setImageResourceData(d, image)
}

func resourceImageRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudcenter.Client)

	image_id_int, err := strconv.Atoi(d.Get("image_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR IMAGE - IMAGE ID INCORRECT")
	}

	image, err := client.GetImage(d.Get("tenant_id").(int), image_id_int)

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR IMAGE")
	}

	return setImageResourceData(d, image)
}

func resourceImageUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	newImage := cloudcenter.Image{

		Id:          d.Get("image_id").(string),
		Name:        d.Get("image_name").(string),
		Description: d.Get("description").(string),
		Visibility:  d.Get("visibility").(string),
		NumOfNICs:   d.Get("num_of_nics").(int),
		OSName:      d.Get("os_name").(string),
		Enabled:     d.Get("enabled").(bool),
		TenantId:    d.Get("tenant_id").(int),
		ImageType:   d.Get("image_type").(string),
	}

	image, err := client.UpdateImage(&newImage)

	if err != nil {
		return errors.New(err.Error())
	}

	return setImageResourceData(d, image)
}

func resourceImageDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*cloudcenter.Client)

	image_id_int, err := strconv.Atoi(d.Get("image_id").(string))

	if err != nil {
		return errors.New("UNABLE TO RETRIEVE DETAILS FOR IMAGE - IMAGE ID INCORRECT")
	}

	err = client.DeleteImage(d.Get("tenant_id").(int), image_id_int)

	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId("")
	return nil
}

func setImageResourceData(d *schema.ResourceData, u *cloudcenter.Image) error {

	if err := d.Set("image_id", u.Id); err != nil {
		return errors.New("CANNOT SET IMAGE ID")
	}
	if err := d.Set("image_name", u.Name); err != nil {
		return errors.New("CANNOT SET IMAGE NAME")
	}
	if err := d.Set("tenant_id", u.TenantId); err != nil {
		return errors.New("CANNOT SET TENANT ID")
	}
	if err := d.Set("resource", u.Resource); err != nil {
		return errors.New("CANNOT SET RESOURCE")
	}
	if err := d.Set("internal_image_name", u.InternalImageName); err != nil {
		return errors.New("CANNOT SET INTERNAL IMAGE NAME")
	}
	if err := d.Set("description", u.Description); err != nil {
		return errors.New("CANNOT SET DESCRIPTION")
	}
	if err := d.Set("visibility", u.Visibility); err != nil {
		return errors.New("CANNOT SET VISIBILITY")
	}
	if err := d.Set("image_type", u.ImageType); err != nil {
		return errors.New("CANNOT SET IMAGE TYPE")
	}
	if err := d.Set("os_name", u.OSName); err != nil {
		return errors.New("CANNOT SET OS NAME")
	}
	/*if err := d.Set("tags", u.Tags); err != nil {
		return errors.New("CANNOT SET TAGS")
	}*/
	if err := d.Set("enabled", u.Enabled); err != nil {
		return errors.New("CANNOT SET ENABLED")
	}
	if err := d.Set("system_image", u.SystemImage); err != nil {
		return errors.New("CANNOT SET SYSTEM IMAGE")
	}
	if err := d.Set("num_of_nics", u.NumOfNICs); err != nil {
		return errors.New("CANNOT SET NUMBER OF NICS")
	}
	if err := d.Set("attach_count", u.AttachCount); err != nil {
		return errors.New("CANNOT SET COUNT")
	}

	return nil
}
