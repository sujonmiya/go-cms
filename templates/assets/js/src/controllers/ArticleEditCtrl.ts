import * as angular from "angular"
import IUploadService = angular.angularFileUpload.IUploadService;
/**
 * Created by Sujon on 12/11/2016.
 */

interface FeaturedImage {
    name: string
    caption: string
    altText: string
}

export class ArticleEditCtrl {
    featuredImage: FeaturedImage;
    file: File;

    static $inject = ['$mdDialog', 'Upload', '$log'];
    constructor(
        private $mdDialog: ng.material.IDialogService,
        private Upload: IUploadService,
        private $log: ng.ILogService){
    }

    upload() {
        this.Upload.upload({
                data: {pic : this.file, featuredImage: this.featuredImage},
                url: '/api/v1/pictures',
                method: 'POST'
            })
            .then((resp) => {
                console.log(resp);
            });
    }

    trashFeaturedImage($event: MouseEvent) {
        this.$log.debug(this.file);

        let confirm = this.$mdDialog.confirm()
            .targetEvent($event)
            .clickOutsideToClose(true)
            .parent(angular.element(document.body))
            .title('Delete?')
            .textContent('Are you sure want to delete the Featured Image "' + this.file + '"?\nThis action can not be undone.')
            .ok('Ok')
            .cancel('Cancel');

        this.$mdDialog.show(confirm)
            .then(function () {
                console.log('Ok');
            }, function () {
                console.log('Cancel');
            })
            .finally(function () {
                confirm = undefined;
            });
    }
}