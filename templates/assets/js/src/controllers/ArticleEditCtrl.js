"use strict";
var ArticleEditCtrl = (function () {
    function ArticleEditCtrl($mdDialog, Upload, $log) {
        this.$mdDialog = $mdDialog;
        this.Upload = Upload;
        this.$log = $log;
    }
    ArticleEditCtrl.prototype.upload = function () {
        this.Upload.upload({
            data: { pic: this.file, featuredImage: this.featuredImage },
            url: '/api/v1/pictures',
            method: 'POST'
        })
            .then(function (resp) {
            console.log(resp);
        });
    };
    ArticleEditCtrl.prototype.trashFeaturedImage = function ($event) {
        this.$log.debug($event.toElement);
        this.$log.debug(this.file);
        var confirm = this.$mdDialog.confirm()
            .targetEvent($event)
            .clickOutsideToClose(true)
            .parent($event.toElement)
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
    };
    ArticleEditCtrl.$inject = ['$mdDialog', 'Upload', '$log'];
    return ArticleEditCtrl;
}());
exports.ArticleEditCtrl = ArticleEditCtrl;
