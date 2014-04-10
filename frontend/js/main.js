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

})

