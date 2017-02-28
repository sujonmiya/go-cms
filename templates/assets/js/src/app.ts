import * as angular from "angular"
//import "jquery"
import "angular-material"
import "angular-animate"
import "angular-sanitize"
import "ng-file-upload"
import "angular-cookies"
import "./node_modules/angular-material-icons/angular-material-icons.min"
import "./node_modules/angular-aria/angular-aria.min"
//import "./node_modules/dropify/dist/js/dropify"
//import "./node_modules/selectize/dist/js/standalone/selectize"
import "./libs/toaster"

import {AppCtrl} from "./controllers/AppCtrl";
import {ArticleCtrl} from "./controllers/ArticleCtrl";
import {ArticleEditCtrl} from "./controllers/ArticleEditCtrl";
import {PagesCtrl} from "./controllers/PagesCtrl";
/**
 * Created by Sujon on 12/11/2016.
 */

module app {
    angular.module('app', [
        'ngMaterial',
        'ngAnimate',
        'ngMdIcons',
        'ngFileUpload',
        'toaster'])
        .config(['$httpProvider', '$interpolateProvider',
            ($httpProvider: ng.IHttpProvider,
             $interpolateProvider: ng.IInterpolateProvider) => {
                $httpProvider.defaults.withCredentials = true;
                $httpProvider.defaults.xsrfCookieName = 'XSRF-TOKEN';
                $httpProvider.defaults.xsrfHeaderName = 'X-XSRF-TOKEN';

                let HEADER_CONTENT_TYPE = 'Content-Type',
                    MEDIA_TYPE_JSON = 'application/json; charset=utf-8';
                $httpProvider.defaults.headers.common.Accept = MEDIA_TYPE_JSON;
                $httpProvider.defaults.headers.common[HEADER_CONTENT_TYPE] = MEDIA_TYPE_JSON;

                $interpolateProvider.startSymbol('[[')
                    .endSymbol(']]');
            }
        ])
        .controller('AppCtrl', AppCtrl)
        .controller('ArticleCtrl', ArticleCtrl)
        .controller('ArticleEditCtrl', ArticleEditCtrl)
        .controller('PagesCtrl', PagesCtrl);
}