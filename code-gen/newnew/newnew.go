package newnew

import (
	"context"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"gitlab.upay.dev/cdfs/txne/constant"
	txneErrors "gitlab.upay.dev/cdfs/txne/errors"
	txneModel "gitlab.upay.dev/cdfs/txne/model"
	"gitlab.upay.dev/cdfs/txne/transform"
	_ "gitlab.upay.dev/cdfs/txnvm/dao"
	txnVMUtility "gitlab.upay.dev/cdfs/txnvm/utility"
	"go.uber.org/zap"
)

type NewNew struct {
	mCM                 model.ConfigMap
	mCapabilityRegistry iface.ICapabilityRegistry
}

func (t *NewNew) Name() string {
	return Name
}

func (t *NewNew) Version() string {
	return constant.Version
}

func (t *NewNew) Category() string {
	return Category
}

func (t *NewNew) ContractId() string {
	return ContractId
}

func (t *NewNew) GetConfigMap() model.ConfigMap {
	return t.mCM
}

func (t *NewNew) SetConfigMap(cm model.ConfigMap) error {
	t.mCM = cm
	return nil
}

func (t *NewNew) SetCapabilityRegistry(capabilityRegistry iface.ICapabilityRegistry) error {
	t.mCapabilityRegistry = capabilityRegistry
	return nil
}

func (t *NewNew) New() iface.ICapability {
	return &NewNew{}
}

func (t *NewNew) Serve(ctx context.Context, event *model.Event) (*model.Event, error) {
	defer func() {
		r := recover()
		if r != nil {
			logger.L(ContractId).Error("panic information", zap.Any("panic", r))
		}
	}()

	payload := &InputModel{}
	if err := transform.UnmarshalJSONEvent(event, payload); err != nil {
		return txneModel.GEE(event.Metadata, t.ContractId(), 400, err), nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return txneModel.GEE(event.Metadata, t.ContractId(), 400, txneErrors.ErrContextDeadlineExceed), nil
	}

	if ctx.Err() == context.Canceled {
		return txneModel.GEE(event.Metadata, t.ContractId(), 400, txneErrors.ErrContextCancelled), nil
	}

	txnVM := txnVMUtility.GetITXNVM(t.mCapabilityRegistry)
	if txnVM == nil {
		return txneModel.GEE(event.Metadata, t.ContractId(), 400, txneErrors.ErrTXNVMNotRegistered), nil
	}

	data, err := txnVM.NewNew(ctx)// Complete method to redirect here
	if err != nil {
		return txneModel.GEE(event.Metadata, t.ContractId(), 400, err), nil
	}

	return txneModel.GSRE(event.Metadata, t.ContractId(), data), nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&NewNew{})
}
