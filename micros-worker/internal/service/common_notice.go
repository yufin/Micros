package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.temporal.io/sdk/client"
	"micros-worker/internal/biz"
	"micros-worker/internal/biz/dto"
	"micros-worker/internal/data"
	"micros-worker/internal/workflow"
	"micros-worker/pkg"
	"net/http"

	pb "micros-worker/api/notice/v1"
)

type CommonNoticeService struct {
	pb.UnimplementedCommonNoticeServer
	log        *log.Helper
	data       *data.Data
	uc         *biz.CommonNoticeUsecase
	cnWorkflow *workflow.CommonNoticeWorkflow
}

func NewCommonNoticeService(logger log.Logger, data *data.Data, cnwf *workflow.CommonNoticeWorkflow, uc *biz.CommonNoticeUsecase) *CommonNoticeService {
	return &CommonNoticeService{
		log:        log.NewHelper(logger),
		data:       data,
		cnWorkflow: cnwf,
		uc:         uc,
	}
}

func (s *CommonNoticeService) SaveSenderConfigWechatBot(ctx context.Context, req *pb.SaveSenderConfigWechatBotReq) (*pb.SaveConfigResp, error) {
	conf := dto.NoticeConfig[dto.WechatBotConfig]{
		Title:        req.Title,
		Comment:      req.Comment,
		SenderConfig: dto.WechatBotConfig{WebHook: req.WebHook},
	}
	conf.Gen()
	b, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}
	docInsertRes, err := s.uc.SaveConfig(ctx, b)
	if err != nil {
		return nil, err
	}
	docId, ok := docInsertRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New(http.StatusFailedDependency, "id assertion failed", "")
	}
	return &pb.SaveConfigResp{
		Success: true,
		Msg:     "success",
		Id:      docId.Hex(),
	}, nil
}

func (s *CommonNoticeService) PubNoticeByWechatBotMarkdown(ctx context.Context, req *pb.PubNoticeByWechatBotMarkdownReq) (*pb.PubNoticeResp, error) {

	wId := fmt.Sprintf("pubNoticeWbMarkdown:%s", pkg.HashStringShort(req.Content+req.SenderId))
	wo := client.StartWorkflowOptions{
		ID:                                       wId,
		TaskQueue:                                "common-notice",
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}

	input := dto.PubNoticeByWechatBotMarkdownReq{
		SenderId: req.SenderId,
		Content:  req.Content,
	}

	_, err := (*s.data.TemporalClient.InfraClient).ExecuteWorkflow(context.Background(), wo, s.cnWorkflow.PubNoticeByWechatBotMarkdown, input)
	if err != nil {
		return nil, err
	}

	return &pb.PubNoticeResp{
		Success:    true,
		Msg:        "success",
		WorkflowId: wId,
	}, nil
}

func (s *CommonNoticeService) PubNoticeByWechatBotText(ctx context.Context, req *pb.PubNoticeByWechatBotTextReq) (*pb.PubNoticeResp, error) {
	return &pb.PubNoticeResp{}, nil
}

//
//func startWorkflowHandler(w http.ResponseWriter, r *http.Request) {
//	// Initialize the Temporal client
//	c, err := client.Dial(client.Options{})
//	if err != nil {
//		// handle error
//	}
//	defer c.Close()
//
//	// Define your workflow options
//	workflowOptions := client.StartWorkflowOptions{
//		// fill in your workflow options
//	}
//
//	// Execute the workflow
//	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, YourWorkflowDefinition, param)
//	if err != nil {
//		// handle error
//	}
//
//	// Send a response back indicating the workflow has started
//	fmt.Fprintf(w, "Workflow started! WorkflowID: %s RunID: %s", workflowRun.GetID(), workflowRun.GetRunID())
//}
