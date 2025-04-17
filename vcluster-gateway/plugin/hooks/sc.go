package hooks

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/loft-sh/vcluster-sdk/hook"
	sdksync "github.com/loft-sh/vcluster-sdk/syncer/context"
	"github.com/pkg/errors"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewScHook(ctx *sdksync.RegisterContext) hook.ClientHook {
	return &scHook{
		rCtx: *ctx,
	}
}

type scHook struct {
	rCtx sdksync.RegisterContext
}

func (s *scHook) Name() string {
	return "sc-hook"
}

func (s *scHook) Resource() client.Object {
	return &storagev1.StorageClass{}
}

var _ hook.MutateCreateVirtual = &scHook{}

func (s *scHook) MutateCreateVirtual(ctx context.Context, obj client.Object) (client.Object, error) {
	sc, ok := obj.(*storagev1.StorageClass)
	if !ok {
		klog.Errorf("Expected a StorageClass object but got: %T", obj)
		return nil, fmt.Errorf("object %v is not a StorageClass", obj)
	}

	//klog.Infof("Starting to mutate StorageClass: %s", sc.Name)

	// Iterate through environment variables and look for matches
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		//klog.V(3).Infof("Processing environment variable: %s", pair[0])

		if !strings.HasPrefix(env, "storage-") {
			//klog.V(4).Infof("Skipping environment variable %s for StorageClass %s", pair[0], sc.Name)
			continue
		}

		if sc.Name == pair[1] {
			klog.Infof("Found matching StorageClass %s for environment variable %s", sc.Name, pair[0])
			return sc, nil
		}
	}

	//klog.Errorf("No matching StorageClass found for any environment variable, StorageClass: %s", sc.Name)
	return nil, errors.New("no matching StorageClass found")
}

var _ hook.MutateUpdateVirtual = &scHook{}

func (s *scHook) MutateUpdateVirtual(ctx context.Context, obj client.Object) (client.Object, error) {
	sc, ok := obj.(*storagev1.StorageClass)
	if !ok {
		klog.Errorf("Expected a StorageClass object but got: %T", obj)
		return nil, fmt.Errorf("object %v is not a StorageClass", obj)
	}

	//klog.Infof("Starting to update StorageClass: %s", sc.Name)

	vsc := &storagev1.StorageClass{}

	// Iterate through environment variables and look for matches
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		//klog.V(3).Infof("Processing environment variable: %s", pair[0])

		if !strings.HasPrefix(env, "storage-") {
			//klog.V(4).Infof("Skipping environment variable %s for StorageClass %s", pair[0], sc.Name)
			continue
		}

		if sc.Name == pair[1] {
			klog.Infof("Found matching StorageClass %s for environment variable %s", sc.Name, pair[0])

			// Check if the virtual StorageClass already exists
			err := s.rCtx.VirtualManager.GetClient().Get(ctx, client.ObjectKey{Name: sc.Name}, vsc)
			if err != nil || vsc.Name == sc.Name {
				klog.Errorf("Failed to retrieve virtual StorageClass %s or StorageClass already exists", sc.Name)
				return nil, errors.New("no matching virtual StorageClass found or it already exists")
			}

			klog.Infof("Successfully mutated and updated StorageClass: %s", sc.Name)
			return sc, nil
		}
	}

	//klog.Errorf("No matching StorageClass found for any environment variable, StorageClass: %s", sc.Name)
	return nil, errors.New("no matching StorageClass found for update")
}

var _ hook.MutateGetPhysical = &scHook{}

func (s *scHook) MutateGetPhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	sc, ok := obj.(*storagev1.StorageClass)
	if !ok {
		klog.Errorf("Expected a StorageClass object but got: %T", obj)
		return nil, fmt.Errorf("object %v is not a StorageClass", obj)
	}

	//klog.Infof("Starting to retrieve physical StorageClass: %s", sc.Name)

	// Retrieve the StorageClass from the physical cluster
	err := s.rCtx.PhysicalManager.GetClient().Get(ctx, client.ObjectKey{Name: sc.Name}, sc)
	if err != nil {
		klog.Errorf("Failed to get physical StorageClass %s: %v", sc.Name, err)
		return nil, fmt.Errorf("failed to get StorageClass %s: %v", sc.Name, err)
	}

	//klog.Infof("Successfully retrieved physical StorageClass: %s", sc.Name)
	return sc, nil
}
