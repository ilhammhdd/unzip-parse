package api

import (
	"bitbucket.com/bbd/unzip-parse/entities"
	"bitbucket.com/bbd/unzip-parse/usecases"
)

func UnzipParse() {
	unzipPool := entities.NewPool(1)
	entities.GoSafely(func() {
		for range unzipPool.Result {
		}
		parsePool := entities.NewPool(1)
		entities.GoSafely(func() {
			for parseResult := range parsePool.Result {
				uploadPool := entities.NewPool(1)
				uploadPool.Work(func() interface{} {
					usecases.UploadPhotos(parseResult.(*[]*usecases.ParseResult))
					return nil
				})
			}
		})
		parsePool.Work(func() interface{} {
			return usecases.Parse()
		})
	})
	unzipPool.Work(func() interface{} {
		usecases.Extractor()
		return nil
	})
}
