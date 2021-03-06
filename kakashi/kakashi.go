package kakashi

import (
	"context"
	"os"

	"github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"go.uber.org/zap"
)

type Kakashi struct {
	logger *zap.Logger

	registry string

	sc *types.SystemContext
	pc *signature.PolicyContext
}

func New(registry, username, password string) (k *Kakashi, err error) {
	logger, _ := zap.NewDevelopment()

	k = &Kakashi{
		logger:   logger,
		registry: registry,
	}

	k.sc = &types.SystemContext{
		// Store authfile in temp, as we will not use it anymore
		AuthFilePath: "/tmp/authfile",
	}

	policy, err := signature.DefaultPolicy(k.sc)
	if err != nil {
		logger.Error("new policy", zap.Error(err))
		return nil, err
	}
	pc, err := signature.NewPolicyContext(policy)
	if err != nil {
		logger.Error("new policy context", zap.Error(err))
		return nil, err
	}
	k.pc = pc

	err = auth.Login(context.Background(), k.sc, &auth.LoginOptions{
		Password: password,
		Username: username,
		Stdout:   os.Stdout,
	}, []string{registry})
	if err != nil {
		k.logger.Error("login registry",
			zap.String("registry", registry),
			zap.Error(err))
		return nil, err
	}

	return
}

func (k *Kakashi) Copy(oldName, newName string) (err error) {
	ctx := context.Background()

	src, err := alltransports.ParseImageName(oldName)
	if err != nil {
		k.logger.Error("parse old image name",
			zap.String("name", oldName),
			zap.Error(err))
		return err
	}
	dst, err := alltransports.ParseImageName(newName)
	if err != nil {
		k.logger.Error("parse new image name",
			zap.String("name", newName),
			zap.Error(err))
		return err
	}

	k.logger.Info("start copy",
		zap.String("src", src.StringWithinTransport()),
		zap.String("dst", dst.StringWithinTransport()))

	_, err = copy.Image(ctx, k.pc, dst, src,
		&copy.Options{
			OptimizeDestinationImageAlreadyExists: true,
			DestinationCtx:                        k.sc,
			ReportWriter:                          os.Stderr,
			ImageListSelection:                    copy.CopyAllImages,
		},
	)
	if err != nil {
		k.logger.Error("copy image", zap.Error(err))
		return err
	}
	k.logger.Info("copied image", zap.String("name", newName))
	return nil
}
