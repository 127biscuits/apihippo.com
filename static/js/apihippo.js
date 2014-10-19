angular.module('hippos', [])
  .controller('HippoController', ['$scope', '$http', function($scope, $http){ 
    $scope.get = function () {
      $http.get('http://random.apihippo.com:8000?verified=false').success(function(data){
        data.url = data.url + "?height=480";
        $scope.hippo = data;
      });
    };

    $scope.vote = function(hippoId) {
      url = 'http://apihippo.com:8000/' + hippoId + '/vote';
      $http.post(url);
      // Don't bother about failures
      $scope.get();
    };

    $scope.hippo = {}
    $scope.get()
  }]);
