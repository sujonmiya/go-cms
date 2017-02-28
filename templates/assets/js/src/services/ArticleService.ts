/**
 * Created by Sujon on 12/15/2016.
 */
import "angular"

interface IArticle {

}

export class ArticleService {

    static $inject = ['$http'];
    constructor(private $http : ng.IHttpService){}

    getArticles(): IArticle[] {
        return
    }

    updateArticle(article: IArticle) {

    }
}