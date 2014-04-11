var chooseKeyword = angular.module(
    'chooseKeyword', 
    [],
    function($interpolateProvider) {
        // Use [[ ]] to delimit AngularJS bindings, because using {{ }} confuses go
        $interpolateProvider.startSymbol('[[');
        $interpolateProvider.endSymbol(']]');
    }
);

chooseKeyword.controller('ChooseKeywordCtrl', function($window, $scope) {
    
    $scope.keywords = $window.keywords;
    $scope.previous = $window.previous;

    $scope.keywordUrl = function (keyword) {
        return '/choose/' + keyword.word + '?previous=' + $scope.previous.join(',');
    }

})

