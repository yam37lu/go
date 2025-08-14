package utils

var (
	ListDistrictVerify              = Rules{"Level": {In("0,1,2,3,4,5,6,7")}, "SubLevel": {In("0,1,2")}, "Extension": {In(",true,false")}}
	AddressVillageQueryVerify       = Rules{"Adcode": {NotSpecialString()}}
	GetTaskProgressVerify           = Rules{"Id": {NotSpecialString()}}
	ListDistrictSubsetVerify        = Rules{"Adcode": {NotSpecialString()}, "Extension": {In(",true,false")}}
	ListUserBindDistrictVerify      = Rules{"Level": {In("1,2,3,4")}}
	ListProcessDetailVerify         = Rules{"ParentAdcode": {NotSpecialString()}, "PageSize": {Le("500"), Ge("1")}, "State": {Le("4"), Ge("0")}}
	ProcessVerify                   = Rules{"Level": {In("0,1,2,3,4,5,6,7")}, "Adcode": {NotSpecialString()}}
	CreateAreaPropertyVerify        = Rules{"Name": {NotEmpty()}, "Type": {NotSpecialString()}}
	EditAreaPropertyVerify          = Rules{"Name": {NotEmpty()}, "Type": {NotSpecialString()}, "Id": {NotSpecialString()}}
	AreaFeatureVerify               = Rules{"Geometry": {NotEmpty()}, "Properties": {NotEmpty()}}
	AreaGeometryVerify              = Rules{"Type": {Eq("Polygon")}, "Coordinates": {NotEmptyDeep()}}
	EditAreaVerify                  = Rules{"Id": {NotSpecialString()}, "Feature": {NotEmpty()}}
	DeleteAreaVerify                = Rules{"Id": {NotSpecialString()}}
	QueryAreaVerify                 = Rules{"Id": {NotSpecialString()}}
	ListAreaVerify                  = Rules{"Bbox": {NotEmpty()}}
	GetAreaTileVerify               = Rules{"Feature": {In("aoi,street,village,county,city,province")}, "Version": {In("v1,v2")}}
	SearchUserVerify                = Rules{"Key": {NotSpecialStringOrEmpty()}}
	BindAndUnbindVerify             = Rules{"Id": {NotEmpty()}, "AddressId": {NotEmpty()}, "DistrictCode": {NotEmpty()}, "Adcode": {NotEmpty()}}
	QueryLogVerify                  = Rules{"Keyword": {NotSpecialStringOrEmpty()}, "PageSize": {Le("500"), Ge("1")}}
	UpdateUserRoleVerify            = Rules{"UserId": {NotSpecialString()}, "RoleCode": {In("sysAdmin,approver,auditor,gridman")}}
	StatisticsProcessListVerify     = Rules{"Level": {In("0,1,2,3,4,5,6,7")}, "Adcode": {NotSpecialString()}, "PageSize": {Le("500"), Ge("1")}}
	StatisticsProcessQueryVerify    = Rules{"Adcode": {NotSpecialString()}, "Level": {In("1,2,3")}}
	QueryBindByAddressDetailVerify  = Rules{"ProvinceName": {NotSpecialString()}, "CityName": {NotSpecialString()}, "CountyName": {NotSpecialString()}}
	SearchShapeInfoVerify           = Rules{"ProvinceName": {NotEmpty()}, "CityName": {NotEmpty()}, "CountyName": {NotEmpty()}}
	ListAddressVerify               = Rules{"Id": {NotSpecialString()}}
	AddConfigVerify                 = Rules{"Key": {NotSpecialString()}, "Value": {NotSpecialString()}, "Type": {Le("100"), Ge("1")}}
	LabelVerify                     = Rules{"Labels": {In("0,1,2,4,8,16,32,64,128,256,512,1024")}}
	AoiSearchByPoiVerify            = Rules{"Id": {NotSpecialString()}}
	RequestWorkflowInfoVerify       = Rules{"Type": {In("1")}, "ApplicantId": {NotEmpty()}, "ApproverId": {NotEmpty()}, "Operation": {In("1,2,3")}, "Aoiid": {NotEmpty()}}
	RequestWorkflowGeomVerify       = Rules{"AoiName": {NotEmpty()}, "AoiType": {NotEmpty()}}
	RequestWorkflowGeomOriginVerify = Rules{"AoiNameOrigin": {NotEmpty()}, "AoiTypeOrigin": {NotEmpty()}}
)
