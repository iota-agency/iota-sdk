package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	model "github.com/iota-agency/iota-sdk/pkg/interfaces/graph/gqlmodels"
)

// UploadFile is the resolver for the uploadFile field.
func (r *mutationResolver) UploadFile(ctx context.Context, file graphql.Upload) (*model.Media, error) {
	//entity := &upload.Upload{
	//	Path:     file.Path,
	//	Path:     file.Path,
	//	Mimetype: file.ContentType,
	//	Size:     file.Size,
	//}
	//
	//if err := r.app.UploadService.UploadFile(ctx, file.File, entity); err != nil {
	//	return nil, err
	//}
	//return entity.ToGraph(), nil
	panic(fmt.Errorf("not implemented: UploadFile - uploadFile"))
}

// DeleteUpload is the resolver for the deleteUpload field.
func (r *mutationResolver) DeleteUpload(ctx context.Context, id int64) (*model.Media, error) {
	panic(fmt.Errorf("not implemented: DeleteUpload - deleteUpload"))
}

// Data is the resolver for the data field.
func (r *paginatedMediaResolver) Data(ctx context.Context, obj *model.PaginatedMedia) ([]*model.Media, error) {
	//entities, err := r.app.UploadService.GetUploadsPaginated(ctx, len(obj.Data), 0, nil)
	//if err != nil {
	//	return nil, err
	//}
	//result := make([]*model.Media, len(entities))
	//for _, entity := range entities {
	//	result = append(result, entity.ToGraph())
	//}
	//return result, nil
	panic(fmt.Errorf("not implemented: Data - data"))
}

// Total is the resolver for the total field.
func (r *paginatedMediaResolver) Total(ctx context.Context, obj *model.PaginatedMedia) (int64, error) {
	//count, err := r.app.UploadService.GetUploadsCount(ctx)
	//if err != nil {
	//	return 0, err
	//}
	//return count, nil
	panic(fmt.Errorf("not implemented: Total - total"))
}

// Upload is the resolver for the upload field.
func (r *queryResolver) Upload(ctx context.Context, id int64) (*model.Media, error) {
	//entity, err := r.app.UploadService.GetUploadByID(ctx, id)
	//if err != nil {
	//	return nil, err
	//}
	//return entity.ToGraph(), nil
	panic(fmt.Errorf("not implemented: Upload - upload"))
}

// Uploads is the resolver for the uploads field.
func (r *queryResolver) Uploads(ctx context.Context, offset int, limit int, sortBy []string) (*model.PaginatedMedia, error) {
	return &model.PaginatedMedia{}, nil
}

// UploadCreated is the resolver for the uploadCreated field.
func (r *subscriptionResolver) UploadCreated(ctx context.Context) (<-chan *model.Media, error) {
	panic(fmt.Errorf("not implemented: UploadCreated - uploadCreated"))
}

// UploadUpdated is the resolver for the uploadUpdated field.
func (r *subscriptionResolver) UploadUpdated(ctx context.Context) (<-chan *model.Media, error) {
	panic(fmt.Errorf("not implemented: UploadUpdated - uploadUpdated"))
}

// UploadDeleted is the resolver for the uploadDeleted field.
func (r *subscriptionResolver) UploadDeleted(ctx context.Context) (<-chan *model.Media, error) {
	panic(fmt.Errorf("not implemented: UploadDeleted - uploadDeleted"))
}

// PaginatedMedia returns PaginatedMediaResolver implementation.
func (r *Resolver) PaginatedMedia() PaginatedMediaResolver { return &paginatedMediaResolver{r} }

type paginatedMediaResolver struct{ *Resolver }
