/**
 * Created by SujonMiya on 28-Feb-16.
 */

(function (angular, undefined) {
    'use strict';

    angular.module('app', ['ngRoute', 'ngMaterial', 'ngAnimate', 'underscore', 'base64', 'toaster'])
        .value('API_URL', '/api/v1')
        .config(['$httpProvider', '$interpolateProvider', function ($httpProvider, $interpolateProvider) {
            $httpProvider.defaults.withCredentials = true;

            var HEADER_CONTENT_TYPE = 'Content-Type',
                MEDIA_TYPE_JSON = 'application/json; charset=utf-8';

            $httpProvider.defaults.headers.common.Accept = MEDIA_TYPE_JSON;
            $httpProvider.defaults.headers.common[HEADER_CONTENT_TYPE] = MEDIA_TYPE_JSON;

            $interpolateProvider.startSymbol('[[').endSymbol(']]');
        }]).directive('ckEditor', function() {
            return {
                require: '?ngModel',
                link: function(scope, elm, attr, ngModel) {
                    if (!ngModel) return;
                    var editor = CKEDITOR.replace(elm[0]);

                    editor.on('pasteState', function() {
                        scope.$apply(function() {
                            ngModel.$setViewValue(editor.getData());
                        });
                    });

                    ngModel.$render = function(value) {
                        editor.setData(ngModel.$viewValue);
                    };
                }
            };
        })
        .factory('ArticleService', ['$http', 'API_URL', function ($http, API_URL) {
            var service = {},
                endpoint = API_URL + '/articles';

            service.getArticles = function () {
                return $http.get(endpoint, {Total: 100, Offset: 0, Sort: '-CreatedAt'});
            };

            service.createArticle = function (article) {
                return $http.post(endpoint, article);
            };

            service.getArticle = function (id) {
                return $http.get(endpoint + '/' +id);
            };

            service.updateArticle = function (article) {
                return $http.put(endpoint + '/'  + article.Id, article);
            };

            service.deleteArticle = function (id) {
                return $http.delete(endpoint + '/'  + id);
            };

            return service;
        }])
        .controller("ArticlesCtrl", ['$scope', 'ArticleService', 'CategoryService', '$mdDialog', 'toaster', function ($scope, ArticleService, CategoryService, $mdDialog, toaster) {
            var event, confirm;
            $scope.articles = [];
            $scope.article = {};
            $scope.categories = [];

            /*ArticleService.getArticles()
                .then(function (response) {
                    $scope.articles = response.data;
                    console.log(response);
                }, function (response) {
                    console.log(response);
                });

            CategoryService.getCategories()
                .then(function (response) {
                    $scope.categories = response.data;
                    console.log(response);
                }, function (response) {
                    console.log(response);
                });*/

            $scope.openMenu = function ($mdOpenMenu, e) {
                event = e;
                $mdOpenMenu(e);
            };

            function remove(article) {
                var index = $scope.articles.indexOf(article);
                $scope.articles.splice(index, 1);
            }

            $scope.delete = function (event, article) {
                confirm = $mdDialog.confirm()
                    .targetEvent(event)
                    .clickOutsideToClose(true)
                    .parent(angular.element(document.body))
                    .title('Delete?')
                    .textContent('Are you sure want to delete the Article "' + article.Title + '"?')
                    .ok('Ok')
                    .cancel('Cancel');

                $mdDialog.show(confirm)
                    .then(function () {
                        console.log('Ok');
                        ArticleService.deleteArticle(article.Id)
                            .then(function (response) {
                                console.log(response);

                                if (response.status === 204) {
                                    remove(article);
                                    toaster.pop('info', 'Hmm...', 'Article ' + article.Title + ' was deleted successfully');
                                }
                            }, function (response) {
                                console.log(response);
                            });
                    }, function () {
                        console.log('Cancel');
                    })
                    .finally(function () {
                        confirm = undefined;
                    });
            };
        }])
        .controller("ArticleCtrl", ['$scope', 'ArticleService', 'CategoryService', 'toaster', '_', function ($scope, ArticleService, CategoryService, toaster, _) {
            $scope.article = {};
            $scope.article.Status = 'Draft';
            $scope.article.Visibility = 'Public';
            $scope.categories = [];

            $scope.getArticle = function(id) {
                ArticleService.getArticle(id)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 200) {
                            $scope.article = response.data;

                            if (!_.isEmpty(response.data.Categories)) {
                                $scope.categories = _.map(response.data.Categories, function (category) {
                                    return category.Id;
                                });
                            }

                            $scope.article.ScheduleAt = new Date(Date.parse(response.data.ScheduleAt));
                        }
                    }, function (response) {
                        console.log(response);
                    });
            };

            $scope.toggle = function (category) {
                var index = _.indexOf($scope.categories, category);
                if (index != -1) {
                    $scope.categories.splice(index, 1);
                }
                else {
                    $scope.categories.push(category);
                }
            };

            $scope.exists = function (category) {
                return _.contains($scope.categories, category);
            };

            $scope.save = function (article) {
                article.Categories = $scope.categories;
                ArticleService.createArticle(article)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 201) {
                            toaster.pop('success', 'Yahoo!', 'Article was created successfully');
                            $scope.article = {};
                            $scope.categories = [];
                        } else {
                            toaster.pop('error', 'Error', response.data.Errors);
                        }
                    }, function (response) {
                        console.log(response);
                        var message = '';
                        angular.forEach(response.data.Errors, function(error) {
                            message += error.FieldName + ' ' + error.Message + '\n';
                        });

                        toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                    });
            };

            $scope.update = function (article) {
                article.Categories = $scope.categories;
                ArticleService.updateArticle(article)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 200) {
                            toaster.pop('info', 'Hmm...', 'Article '+ article.Title + ' was updated successfully');
                        }
                    }, function (response) {
                        console.log(response);
                        var message = '';
                        angular.forEach(response.data.Errors, function(error) {
                            message += error.FieldName + ' ' + error.Message + '\n';
                        });

                        toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                    });
            };
        }])
        .factory('PageService', ['$http', 'API_URL', function ($http, API_URL) {
            var service = {},
                endpoint = API_URL + '/pages';

            service.getPages = function () {
                return $http.get(endpoint, {Total: 100, Offset: 0, Sort: '-CreatedAt'});
            };

            service.createPage = function (page) {
                return $http.post(endpoint, page);
            };

            service.getPage = function (id) {
                return $http.get(endpoint + '/' +id);
            };

            service.updatePage = function (page) {
                return $http.put(endpoint + '/'  + page.Id, page);
            };

            service.deletePage = function (id) {
                return $http.delete(endpoint + '/'  + id);
            };

            return service;
        }])
        .controller("PagesCtrl", ['$scope', 'PageService', '$mdDialog', 'toaster', '_', function ($scope, PageService, $mdDialog, toaster, _) {
            var event, confirm;
            $scope.pages = [];
            $scope.page = {};
            $scope.filter = {};
            $scope.authors = [];
            $scope.editors = [];
            $scope.statuses = [];
            $scope.visibilities = [];

            PageService.getPages()
                .then(function (response) {
                    $scope.pages = response.data;
                    console.log(response);

                    _.each(angular.copy($scope.pages), function (page) {
                        if (!_.contains($scope.statuses, page.Status))
                            $scope.statuses.push(page.Status);

                        if (!_.contains($scope.visibilities, page.Visibility))
                            $scope.visibilities.push(page.Visibility);

                        if (!_.contains($scope.authors, page.Author.FullName))
                            $scope.authors.push(page.Author.FullName);

                        if(page.Editor)
                            if (!_.contains($scope.editors, page.Editor.FullName))
                                $scope.editors.push(page.Editor.FullName);
                    })
                }, function (response) {
                    console.log(response);
                });

            $scope.openMenu = function ($mdOpenMenu, e) {
                event = e;
                $mdOpenMenu(e);
            };

            function remove(page) {
                var index = $scope.pages.indexOf(page);
                $scope.pages.splice(index, 1);
            }

            $scope.delete = function (event, page) {
                confirm = $mdDialog.confirm()
                    .targetEvent(event)
                    .clickOutsideToClose(true)
                    .parent(angular.element(document.body))
                    .title('Delete?')
                    .textContent('Are you sure want to delete the Page "' + page.Title + '"?')
                    .ok('Ok')
                    .cancel('Cancel');

                $mdDialog.show(confirm)
                    .then(function () {
                        console.log('Ok');
                        PageService.deletePage(page.Id)
                            .then(function (response) {
                                console.log(response);

                                if (response.status === 204) {
                                    remove(page);
                                    toaster.pop('info', 'Hmm...', 'Page ' + page.Title + ' was deleted successfully');
                                }
                            }, function (response) {
                                console.log(response);
                            });
                    }, function () {
                        console.log('Cancel');
                    })
                    .finally(function () {
                        confirm = undefined;
                    });
            };
        }])
        .controller("PageCtrl", ['$scope', 'PageService', 'toaster', '_', function ($scope, PageService, toaster, _) {
            $scope.page = {};
            $scope.page.Template = 'Default';
            $scope.page.Status = 'Draft';
            $scope.page.Visibility = 'Public';

            $scope.getPage = function(id) {
                PageService.getPage(id)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 200) {
                            $scope.page = response.data;
                            $scope.page.ScheduleAt = new Date(Date.parse(response.data.ScheduleAt));
                        }
                    }, function (response) {
                        console.log(response);
                    });
            };

            $scope.save = function (page) {
                PageService.createPage(page)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 201) {
                            $scope.page = {};
                            toaster.pop('success', 'Yahoo!', 'Page was created successfully');
                        } else {
                            toaster.pop('error', 'Error', response.data.Errors);
                        }
                    }, function (response) {
                        console.log(response);
                        var message = '';
                        angular.forEach(response.data.Errors, function(error) {
                            message += error.FieldName + ' ' + error.Message + '\n';
                        });

                        toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                    });
            };

            $scope.update = function (page) {
                PageService.updatePage(page)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 200) {
                            toaster.pop('info', 'Hmm...', 'Page '+ page.Title + ' was updated successfully');
                        }
                    }, function (response) {
                        console.log(response);
                        var message = '';
                        angular.forEach(response.data.Errors, function(error) {
                            message += error.FieldName + ' ' + error.Message + '\n';
                        });

                        toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                    });
            };
        }])
        .factory('UserService', ['$http', 'API_URL', function ($http, API_URL) {
            var service = {},
                endpoint = API_URL + '/users';

            service.getUsers = function () {
                return $http.get(endpoint, {Total: 100, Offset: 0, Sort: '-CreatedAt'});
            };

            service.createUser = function (user) {
                return $http.post(endpoint, user);
            };

            service.updateUser = function (user) {
                return $http.put(endpoint + '/' + user.Id, user);
            };

            service.deleteUser = function (id) {
                return $http.delete(endpoint + '/'  + id);
            };

            service.getRoles = function () {
                return $http.get(endpoint + '/roles');
            };

            return service;
        }])
        .controller("UsersCtrl", ['$scope', 'UserService', '$mdDialog', '$base64', 'toaster', function ($scope, UserService, $mdDialog, $base64, toaster) {
            var event, confirm;
            $scope.users = [];
            $scope.roles = [];

            /*UserService.getUsers()
                .then(function (response) {
                    $scope.users = response.data;
                    console.log(response);
                }, function (response) {
                    console.log(response);
                });

            UserService.getRoles()
                .then(function (response) {
                    $scope.roles = response.data;
                    console.log(response);
                }, function (response) {
                    console.log(response);
                });*/

            $scope.openMenu = function ($mdOpenMenu, e) {
                event = e;
                $mdOpenMenu(e);
            };

            function remove(user) {
                var index = $scope.users.indexOf(user);
                $scope.users.splice(index, 1);
            }

            $scope.delete = function (event, user) {
                confirm = $mdDialog.confirm()
                    .targetEvent(event)
                    .clickOutsideToClose(true)
                    .parent(angular.element(document.body))
                    .title('Delete?')
                    .textContent('Are you sure want to delete the User "' + user.FirstName + ' ' + user.LastName + '"?')
                    .ok('Ok')
                    .cancel('Cancel');

                $mdDialog.show(confirm)
                    .then(function () {
                        console.log('Ok');
                        UserService.deleteUser(user.Id)
                            .then(function (response) {
                                console.log(response);

                                if (response.status === 204) {
                                    console.log($scope);
                                    remove(user);
                                    toaster.pop('info', 'Hmm...', 'User ' + user.FirstName + ' ' + user.LastName + ' was deleted successfully');
                                }
                            }, function (response) {
                                console.log(response);
                            });
                    }, function () {
                        console.log('Cancel');
                    })
                    .finally(function () {
                        confirm = undefined;
                    });
            };

            function DialogController($scope, $mdDialog) {
                $scope.hide = function () {
                    $mdDialog.hide();
                };

                $scope.cancel = function () {
                    $mdDialog.cancel();
                };
            }

            $scope.save = function (user) {
                user.Password = $base64.encode(user.Password);
                UserService.createUser(user)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 201) {
                            $scope.users.push(response.data);
                            toaster.pop('success', 'Yahoo!', 'User ' + user.FirstName + ' ' + user.LastName + ' was created successfully');
                            $scope.user = {};
                            $mdDialog.hide();
                        }
                    }, function (response) {
                        console.log(response);
                        var message = '';
                        angular.forEach(response.data.Errors, function(error) {
                            message += error.FieldName + ' ' + error.Message + '\n';
                        });

                        toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                    });
            };

            $scope.newUser = function ($event) {
                $scope.user = undefined;
                var options = {
                    controller: DialogController,
                    templateUrl: '/assets/dialogs/user.html',
                    parent: angular.element(document.body),
                    targetEvent: $event,
                    clickOutsideToClose: false,
                    scope: $scope,
                    preserveScope: true
                };

                $mdDialog.show(options);
            };

            function update(user) {
                var index = _.findIndex($scope.users, function (u) {
                    return _.isEqual(u.Id, user.Id);
                });

                $scope.users[index] = angular.copy(user);
            }

            $scope.update = function (user) {
                if (user.Password) {
                    user.Password = $base64.encode(user.Password);
                }

                UserService.updateUser(user)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 200) {
                            update(user);
                            toaster.pop('info', 'Hmm...', 'User ' + user.FirstName + user.LastName + ' was updated successfully');
                            $mdDialog.hide();
                        }
                    }, function (response) {
                        console.log(response);
                        var message = '';
                        angular.forEach(response.data.Errors, function(error) {
                            message += error.FieldName + ' ' + error.Message + '\n';
                        });

                        toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                    });
            };

            $scope.edit = function (event, user) {
                console.log(user);
                $scope.user = angular.copy(user);

                var options = {
                    controller: DialogController,
                    templateUrl: '/assets/dialogs/user-edit.html',
                    parent: angular.element(document.body),
                    targetEvent: event,
                    clickOutsideToClose: false,
                    scope: $scope,
                    preserveScope: true
                };

                $mdDialog.show(options);
            };
        }])
        .factory('CategoryService', ['$http', 'API_URL', function ($http, API_URL) {
            var service = {},
                endpoint = API_URL + '/categories';

            service.getCategories = function () {
                return $http.get(endpoint, {Total: 100, Offset: 0, Sort: '-CreatedAt'});
            };

            service.createCategory = function (category) {
                return $http.post(endpoint, category);
            };

            service.updateCategory = function (category) {
                return $http.put(endpoint + category.Id, category);
            };

            service.deleteCategory = function (id) {
                return $http.delete(endpoint + id);
            };

            return service;
        }])
        .controller("CategoriesCtrl", ['$scope', 'CategoryService', '$mdDialog', '_', 'toaster', function ($scope, CategoryService, $mdDialog, _, toaster) {
            var event, confirm;
            $scope.categories = [];
            $scope.category = {};

            /*CategoryService.getCategories()
                .then(function (response) {
                    $scope.categories = response.data;
                    console.log(response);
                }, function (response) {
                    console.log(response);
                    var message = '';
                    angular.forEach(response.data.Errors, function(error) {
                        message += error.FieldName + ' ' + error.Message + '\n';
                    });

                    toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                });*/

            $scope.openMenu = function ($mdOpenMenu, e) {
                event = e;
                $mdOpenMenu(e);
            };

            function remove(category) {
                var index = $scope.categories.indexOf(category);
                $scope.categories.splice(index, 1);
            }

            $scope.delete = function (event, category) {
                confirm = $mdDialog.confirm()
                    .targetEvent(event)
                    .clickOutsideToClose(true)
                    .parent(angular.element(document.body))
                    .title('Delete?')
                    .textContent('Are you sure want to delete the Category "' + category.Name + '"?')
                    .ok('Ok')
                    .cancel('Cancel');

                $mdDialog.show(confirm)
                    .then(function () {
                        console.log('Ok');
                        CategoryService.deleteCategory(category.Id)
                            .then(function (response) {
                                console.log(response);

                                if (response.status === 204) {
                                    console.log($scope);
                                    remove(category);
                                    toaster.pop('info', 'Hmm...', 'Category ' + category.Name + ' was deleted successfully');
                                }
                            }, function (response) {
                                console.log(response);
                            });
                    }, function () {
                        console.log('Cancel');
                    })
                    .finally(function () {
                        confirm = undefined;
                    });
            };

            function DialogController($scope, $mdDialog) {
                $scope.hide = function () {
                    $mdDialog.hide();
                };

                $scope.cancel = function () {
                    $mdDialog.cancel();
                };
            }

            $scope.save = function (category) {
                CategoryService.createCategory(category)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 201) {
                            $scope.categories.push(response.data);
                            $scope.category = {};
                            toaster.pop('success', 'Yahoo!', 'Category ' + response.data.Name + ' was created successfully');
                        } else {
                            toaster.pop('error', 'Error', response.data.Errors);
                        }
                    }, function (response) {
                        console.log(response);
                        var message = '';
                        angular.forEach(response.data.Errors, function(error) {
                            message += error.FieldName + ' ' + error.Message + '\n';
                        });

                        toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                    });
            };

            $scope.newCategory = function ($event) {
                var options = {
                    controller: DialogController,
                    templateUrl: '/assets/dialogs/category.html',
                    parent: angular.element(document.body),
                    targetEvent: $event,
                    clickOutsideToClose: false,
                    scope: $scope,
                    preserveScope: true
                };

                $mdDialog.show(options);
            };

            function update(category) {
                var index = _.findIndex($scope.categories, function (c) {
                    return _.isEqual(c.Id, category.Id);
                });

                $scope.categories[index] = angular.copy(category);
            }

            $scope.update = function (category) {
                console.log(category);
                console.log(category.Parent);
                if(category.Parent) {
                    category.Parent = category.Parent.Id;
                }

                CategoryService.updateCategory(category)
                    .then(function (response) {
                        console.log(response);

                        if (response.status === 200) {
                            update(category);
                            toaster.pop('info', 'Hmm...', 'Category ' + category.Name + ' was updated successfully');
                            $mdDialog.hide();
                        }
                    }, function (response) {
                        console.log(response);
                        var message = '';
                        angular.forEach(response.data.Errors, function(error) {
                            message += error.FieldName + ' ' + error.Message + '\n';
                        });

                        toaster.pop('error', 'Oops, ' + response.data.Reason, message);
                    });
            };

            $scope.edit = function (event, category) {
                console.log(category);
                $scope.category = angular.copy(category);
                $scope.parents = _.without($scope.categories, category);

                var options = {
                    controller: DialogController,
                    templateUrl: '/assets/dialogs/category-edit.html',
                    parent: angular.element(document.body),
                    targetEvent: event,
                    clickOutsideToClose: false,
                    scope: $scope,
                    preserveScope: true
                };

                $mdDialog.show(options);
            };
        }])
        .controller('PicturesCtrl', function () {

        });

})(angular);