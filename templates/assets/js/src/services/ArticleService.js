"use strict";
/**
 * Created by Sujon on 12/15/2016.
 */
require("angular");
var ArticleService = (function () {
    function ArticleService($http) {
        this.$http = $http;
    }
    ArticleService.prototype.getArticles = function () {
        return;
    };
    ArticleService.prototype.updateArticle = function (article) {
    };
    ArticleService.$inject = ['$http'];
    return ArticleService;
}());
exports.ArticleService = ArticleService;
