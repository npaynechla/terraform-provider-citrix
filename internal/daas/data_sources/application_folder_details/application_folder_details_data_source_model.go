// Copyright © 2023. Citrix Systems, Inc.

package application_folder_details

import (
	"github.com/citrix/citrix-daas-rest-go/citrixorchestration"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApplicationFolderDetailsDataSourceModel struct {
	Path              types.String               `tfsdk:"path"`
	TotalApplications types.Int64                `tfsdk:"total_applications"`
	ApplicationsList  []ApplicationResourceModel `tfsdk:"applications_list"`
}

type ApplicationResourceModel struct {
	Name                   types.String      `tfsdk:"name"`
	PublishedName          types.String      `tfsdk:"published_name"`
	Description            types.String      `tfsdk:"description"`
	InstalledAppProperties InstalledAppModel `tfsdk:"installed_app_properties"`
	DeliveryGroups         []types.String    `tfsdk:"delivery_groups"`
	ApplicationFolderPath  types.String      `tfsdk:"application_folder_path"`
}

type InstalledAppModel struct {
	CommandLineArguments  types.String `tfsdk:"command_line_arguments"`
	CommandLineExecutable types.String `tfsdk:"command_line_executable"`
	WorkingDirectory      types.String `tfsdk:"working_directory"`
}

func (r ApplicationFolderDetailsDataSourceModel) RefreshPropertyValues(apps *citrixorchestration.ApplicationResponseModelCollection) ApplicationFolderDetailsDataSourceModel {

	var res []ApplicationResourceModel
	for _, app := range apps.GetItems() {
		res = append(res, ApplicationResourceModel{
			Name:                   types.StringValue(app.GetName()),
			PublishedName:          types.StringValue(app.GetPublishedName()),
			Description:            types.StringValue(app.GetDescription()),
			ApplicationFolderPath:  types.StringValue(*app.GetApplicationFolder().Name.Get()),
			InstalledAppProperties: r.getInstalledAppProperties(app),
			DeliveryGroups:         r.getDeliveryGroups(app),
		})
	}

	r.ApplicationsList = res
	r.TotalApplications = types.Int64Value(int64(*apps.TotalItems.Get()))
	return r
}

func (r ApplicationFolderDetailsDataSourceModel) getInstalledAppProperties(app citrixorchestration.ApplicationResponseModel) InstalledAppModel {
	return InstalledAppModel{
		CommandLineArguments:  types.StringValue(app.GetInstalledAppProperties().CommandLineArguments),
		CommandLineExecutable: types.StringValue(app.GetInstalledAppProperties().CommandLineExecutable),
		WorkingDirectory:      types.StringValue(app.GetInstalledAppProperties().WorkingDirectory),
	}
}

func (r ApplicationFolderDetailsDataSourceModel) getDeliveryGroups(app citrixorchestration.ApplicationResponseModel) []types.String {
	var res []types.String
	for _, dg := range app.AssociatedDeliveryGroupUuids {
		res = append(res, types.StringValue(dg))
	}
	return res
}
