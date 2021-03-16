package services

//
//import (
//	"github.com/phachon/mm-wiki/app/models"
//	"github.com/phachon/mm-wiki/app/utils"
//)
//
//type Doc struct {
//}
//
//type Error struct {
//	errorLog  string
//	errorInfo string
//}
//
//var DocService = newDocService()
//
//func newDocService() *Doc {
//	return &Doc{
//	}
//}
//
//func (ds *Doc) GetDocumentPageData(documentId string) ( Data map[string]string , docError Error){
//
//	if documentId == "" {
//		docError.errorInfo = "文档未找到"
//		return
//	}
//
//	document, err := models.DocumentModel.GetDocumentByDocumentId(documentId)
//	if err != nil {
//		docError.errorLog = "修改文档 " + documentId + " 失败：" + err.Error()
//		docError.errorInfo = "修改文档失败！";
//		return
//	}
//	if len(document) == 0 {
//		docError.errorInfo = "文档不存在！";
//		return
//	}
//
//	spaceId := document["space_id"]
//	space, err := models.SpaceModel.GetSpaceBySpaceId(spaceId)
//	if err != nil {
//		docError.errorLog = "修改文档 " + documentId + " 失败：" + err.Error()
//		docError.errorInfo = "修改文档失败！";
//		return
//	}
//	if len(space) == 0 {
//		docError.errorInfo = "文档所在空间不存在！";
//		return
//	}
//	// check space visit_level
//	_, isEditor, _ := this.GetDocumentPrivilege(space)
//	if !isEditor {
//		this.ViewError("您没有权限修改该空间下文档！")
//	}
//
//	// get parent documents by document
//	_, pageFile, err := models.DocumentModel.GetParentDocumentsByDocument(document)
//	if err != nil {
//		this.ErrorLog("查找父文档失败：" + err.Error())
//		this.ViewError("查找父文档失败！")
//	}
//
//	// get document content
//	documentContent, err := utils.Document.GetContentByPageFile(pageFile)
//	if err != nil {
//		this.ErrorLog("查找文档 " + documentId + " 失败：" + err.Error())
//		this.ViewError("文档不存在！")
//	}
//
//	autoFollowDoc := models.ConfigModel.GetConfigValueByKey(models.ConfigKeyAutoFollowdoc, "0")
//	sendEmail := models.ConfigModel.GetConfigValueByKey(models.ConfigKeySendEmail, "0")
//}
//
//func (this *Doc) GetDocumentPrivilege(space map[string]string) (isVisit, isEditor, isManager bool) {
//
//	if this.IsRoot() {
//		return true, true, true
//	}
//	spaceUser, _ := models.SpaceUserModel.GetSpaceUserBySpaceIdAndUserId(space["space_id"], this.UserId)
//	if len(spaceUser) == 0 {
//		if space["visit_level"] == models.Space_VisitLevel_Private {
//			return false, false, false
//		} else {
//			return true, false, false
//		}
//	}
//	if spaceUser["privilege"] == fmt.Sprintf("%d", models.SpaceUser_Privilege_Editor) {
//		return true, true, false
//	}
//	if spaceUser["privilege"] == fmt.Sprintf("%d", models.SpaceUser_Privilege_Manager) {
//		return true, true, true
//	}
//	return true, false, false
//}
