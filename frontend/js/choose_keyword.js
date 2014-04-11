var sourgrapes = angular.module(
    'sourgrapes', 
    [],
    function($interpolateProvider) {
        // Use [[ ]] to delimit AngularJS bindings, because using {{ }} confuses go
        $interpolateProvider.startSymbol('[[');
        $interpolateProvider.endSymbol(']]');
    }
);

sourgrapes.controller('ChooseKeywordCtrl', function($window, $scope) {
    
    $scope.keywords = $window.keywords;
    $scope.previous = $window.previous;

    $scope.keywordUrl = function (keyword) {
        return '/keyword/' + keyword.word + '?previous=' + $scope.previous.join(',');
    }

    $scope.imageUrl = function(keyword) {
        return 'https://www.google.com/search?q=' + keyword.word + '&tbm=isch';
    }

})

