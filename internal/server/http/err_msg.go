package http

// err_code for response

/**
上传参数相关错误
 */
const ERR_INPUT_NULL = -101			// 缺少参数
const ERR_INPUT_TYPE = -102			//参数类型错误
const ERR_INPUT_NO_PHOTO = -103		//上传照片失败
const ERR_INPUT_NO_FILE = -104		//上传文件失败
const ERR_NOT_ACCORD = -105			//不符合的参数

/**
koala相关错误
 */
const ERR_KOALA_DATA = -201			//获取考拉信息失败
const ERR_KOALA_DO = -202				//操作考拉接口失败
const ERR_KOALA_LOGIN = -204			//考拉登陆失败

/**
命令行调用
 */
const ERR_COMMAND  = -300				//命令行调用失败

/**
SQL相关错误
 */
const ERR_SQL_DATE = -402				//数据库操作失败
const ERR_GET_INFO = -404				//获取信息失败

/**
系统/未知错误
 */
const ERR_SYSTEM = -500				//系统错误
const ERR_INPUT_PARAMETER = -501		//参数错误

/**
登陆统一提示信息
 */
const ERR_LOGIN = -401				//登陆失败
const ERR_LOGIN_AUTHORITY = -100		//无权限访问


/**
其他错误
 */
const ERR_COPY_PHOTO  = -600			//copy图片失败
const ERR_WRITE_YAML  = -601			//yaml写入失败
const ERR_JSON_CONVERSION  = -602	//转化json出错
const ERR_DIR_CREATE  = -603			//创建文件夹出错
const ERR_FILE_CREATE  = -604		//创建文件出错
