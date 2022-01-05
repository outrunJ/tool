# encoding: utf8
import jinja2

tpl_text = '''package {{ package }}

import "git.meiqia.com/business_platform/component/response"

var (
    // res
	ResNewResponseWithCode       func(int, string, string, interface{}) *response.Response
    {%- for code, msg in codemsg.items() %}
    Res{{msg|replace("_", "")}} func(string) *response.Response
    {%- endfor %}
)

func InitRes(res response.Responder) {
	if res == nil {
    	return 
    }
    
    var errCodeMap = map[int]string{
    {%- for code, msg in codemsg.items() %}
    	{{code}}: "Territory_{{msg}}",
    {%- endfor %}
    }
    
    ResNewResponseWithCode = func(code int, msg string, desc string, stack interface{}) *response.Response {
		return &response.Response{
			Code:    code,
			MessageCommon: msg,
			Message: desc,
			Stack: stack,
		}
	}
    
    resDiy := res.Map(errCodeMap)
    {%- for code, msg in codemsg.items() %}
    	Res{{msg|replace("_","")}} = func(msg string) *response.Response {
    		return resDiy({{code}},msg)
    	}
    {%- endfor %}
}

'''

codemsg = {
    101: "No_Permission",
    102: "Type_Name_Already_Exists",
    103: "Claim_Limit",
    104: "Model_Limit",
    105: "Running_Model_Already_Exists",
    106: "Tenant_Admin_Not_Exists",
    107: "Parameter_Empty",
    108: "No_Running_Model",
    109: "Record_Already_Exists",
    110: "DB_Insert_Error",
    111: "Select_Field_Name_Unreadable",
    112: "DB_Update_Error",
    113: "DB_Delete_Error",
    114: "Territory_Capacity_Limit",
    115: "DB_Select_Error",
    116: "Territory_Record_Not_Empty",
    117: "Territory_Property_Incompatible",
    118: "Move_Restrict_Passive",
    120: "Invalid_Type_Name",
    121: "Invalid_Model_Id",
    122: "Invalid_Territory_Id",
    123: "Invalid_Record_Id",
    124: "Invalid_User_Id",
    125: "Invalid_Member_Id",
    126: "Invalid_Field_Name",
    127: "Invalid_Permission_Id",
    128: "Invalid_Token",
    130: "Service_Judge_Not_Found",
    131: "Service_Mysql_Not_Found",
    132: "Service_Tenant_Not_Found",
    133: "Service_Idg_Not_Found",
    134: "Service_MDS_DDL_Not_Found",
    135: "Service_MQ_Not_Found",
    136: "Service_MQ_Serialize_Error",
    137: "Service_MQ_Deserialize_Error",
    138: "Service_MQ_Publish_Error",
    139: "Service_MQ_Message_Not_Found",
    140: "Service_MDS_Not_Found",
    200: "Property_Not_Created",
    201: "Property_Not_Deleted",
    202: "Property_Not_Updated",
    203: "Property_Not_Found",
    204: "Property_Invalid",
    205: "Property_Not_Be_Given",
    206: "Property_Not_Be_Taken",
    207: "Script_Not_Created",
    208: "Script_Not_Deleted",
    209: "Script_Not_Updated",
    210: "Script_Not_Found",
    211: "Procedure_Not_Created",
    212: "Procedure_Not_Deleted",
    213: "Procedure_Not_Updated",
    214: "Procedure_Not_Found",
    215: "Procedure_Invalid",
    216: "Procedure_Formula_Invalid",
    223: "Procedure_Params_Invalid",
    217: "Territory_Procedure_Not_Created",
    218: "Territory_Procedure_Not_Deleted",
    219: "Territory_Procedure_Not_Updated",
    220: "Territory_Procedure_Not_Found",
    221: "Motion_Invalid",
    222: "Territory_Not_Found",
    230: "Record_Not_Found",
    240: "Service_IDG_Timeout",
    241: "Service_Judge_Timeout",
    242: "Service_MDS_Timeout",
    243: "Service_Tenant_Timeout",
    250: "Type_Not_Found",
    900: "Internal_Error",
}


def generate():
    tpl = jinja2.Template(tpl_text)

    return tpl.render(**{
        'package': 'erro',
        'codemsg': codemsg,
    })


def main():
    print(generate())


if __name__ == '__main__':
    main()
