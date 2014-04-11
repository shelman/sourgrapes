var movieSearch = angular.module(
    'movieSearch', 
    [],
    function($interpolateProvider) {
        // Use [[ ]] to delimit AngularJS bindings, because using {{ }} confuses go
        $interpolateProvider.startSymbol('[[');
        $interpolateProvider.endSymbol(']]');
    }
);

movieSearch.controller('MovieSearchCtrl', function($scope, $http) {

    $scope.searchVal = '';
    $scope.searchResults = [];
    $scope.$watch('searchVal', function(newVal, oldVal) {
        if (newVal === '') {
            $scope.searchResults = [];
            return;
        }
        $http.get('/search_results/movie/' + newVal).
            success(function(data, status) {
                $scope.searchResults = data.results;
            });

    });

})
