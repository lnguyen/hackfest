var app = angular.module('hackfest', []);

app.controller('CoinCtrl', function($scope, $http, $timeout) {
  (function tick() {
    $http.get('coin').success(function (data) {
      console.log(data);
      if (data["coin"] == "Found") {
        $scope.coinEject = false;
        $scope.coinFound = true;
        $timeout(tick, 2000);
      } else if (data["coin"] == "Ejected") {
        $scope.coinFound = false;
        $scope.coinEject = true;
        $timeout(tick, 2000);
      } else {
        $scope.coinFound = false;
        $scope.coinEject = false;
        $timeout(tick, 1000);
      }
    });
  })();
});
