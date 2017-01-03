/***
 * angularjs ʹ��ָ���װ�ٶȵ�ͼ
 * Created by xc on 2015/6/13.
 */
(function (window, angular) {
  'use strict';
  var angularMapModule = angular.module('angularBMap', []);
  angularMapModule.provider('angularBMap', mapService);//��λ����
  angularMapModule.directive('angularBmap', mapDirective);//��λָ��
  /*
   * ��λ��ط���
   */
  function mapService() {
    //��������
    this.default_position = new BMap.Point(118.789572, 32.048667);//��ͼĬ�����ĵ�
    /**
     * ���õ�ͼĬ�����ĵ�
     * @param lng
     * @param lat
     * @returns {mapService}
     */
    this.setDefaultPosition = function (lng, lat) {
      this.default_position = new BMap.Point(lng, lat);
      return this;
    };

    //���صķ���
    this.$get = BMapService;
    BMapService.$inject = ['$q', '$timeout'];
    function BMapService($q, $timeout) {
      var map,//ȫ�ֿ��õ�map����
        default_position = this.default_position;//Ĭ�����ĵ�
      return {
        initMap: initMap,//��ʼ����ͼ
        getMap: getMap,//���ص�ǰ��ͼ����
        geoLocation: geoLocation,//��ȡ��ǰλ��
        geoLocationAndCenter: geoLocationAndCenter,//��ȡ��ǰλ�ã�������ͼ�ƶ�����ǰλ��
        drawMarkers: drawMarkers,//�����Ȥ��
        drawMarkersAndSetOnclick: drawMarkersAndSetOnclick//�����Ȥ��ͬʱ��ӵ���¼�
      };
      /**
       * ��ȡmap����
       * @alias getMap
       */
      function getMap() {
        if (!map) {
          map = new BMap.Map('bMap');//��ͼ����
        }
        return map;
      }

      /**
       * ��ʼ����ͼ
       * @constructor
       */
      function initMap() {
        var defer = $q.defer();
        $timeout(function () {
          getMap().centerAndZoom(default_position, 14);
          defer.resolve();
        });
        return defer.promise;
      }

      /**
       * ���ðٶȵ�ͼ��ȡ��ǰλ��
       * @constructor
       */
      function geoLocation() {
        var defer = $q.defer(), location = new BMap.Geolocation();//�ٶȵ�ͼ��λʵ��
        location.getCurrentPosition(function (result) {
          if (this.getStatus() === BMAP_STATUS_SUCCESS) {
            //��λ�ɹ�,���ض�λ�ص�;���
            defer.resolve(result);
          } else {
            defer.reject('���ܻ�ȡλ��');
          }
        }, function (err) {
          defer.reject('��λʧ��');
        });
        return defer.promise;
      }

      /**
       * ��ȡ��ǰλ�ã�������ͼ�ƶ�����ǰλ��
       * @constructor
       */
      function geoLocationAndCenter() {
        var defer = $q.defer();
        geoLocation().then(function (result) {
          getMap().panTo(result.point);
          var marker = new BMap.Marker(result.point);
          getMap().addOverlay(marker);
          defer.resolve(result);
        }, function (err) {
          //��λʧ��
          getMap().panTo(default_position);
          var marker = new BMap.Marker(default_position);
          getMap().addOverlay(marker);
          defer.reject('��λʧ��');
        });
        return defer.promise;
      }

      /**
       * ���ͼ�����Ȥ�㣨marker��
       * @param markers
       */
      function drawMarkers(markers) {
        var _markers = [],//����ӵ���Ȥ���б�
          defer = $q.defer(),
          point,//��ǰ��ӵ������
          _length,//���鳤��
          _progress;//��ǰ������ӵĵ������
        $timeout(function () {
          //�ж��Ƿ��ж�λ��
          if (!markers) {
            defer.reject('û�д�����Ȥ��');
            return;
          }
          //�����˲���
          if (!angular.isArray(markers)) {
            //����Ĳ���array
            if (markers.loc) {
              _markers.push(markers);
            } else {
              defer.reject('��ȡ����loc������Ϣ');
            }
          } else {
            if (markers[0].loc) {
              _markers = markers;
            } else {
              defer.reject('��ȡ����loc������Ϣ');
            }
          }
          _length = _markers.length - 1;
          angular.forEach(_markers, function (obj, index) {
            _progress = index;
            if (angular.isObject(obj.loc)) {
              point = new BMap.Point(obj.loc.lng, obj.loc.lat);
            } else if (angular.isString(obj.loc)) {
              point = new BMap.Point(obj.loc.split(',')[0], obj.loc.split(',')[1]);
            } else {
              _progress = '��' + index + '����Ȥ��loc���󲻴��ڻ��ʽ����ֻ֧��object��string';
            }
            var marker = new BMap.Marker(point);
            getMap().addOverlay(marker);
            defer.notify(_progress);
            if (index === _length) {
              defer.resolve();
            }
          });
        });
        return defer.promise;
      }

      /**
       * Ĭ�ϵ���¼�
       * @param obj
       */
      function markerClick() {
        getMap().panTo(this.point);
      }

      /**
       * ���ͼ�����Ȥ��ͬʱ��ӵ���¼�
       * @param markers
       * @param onClick
       * @returns {*}
       */
      function drawMarkersAndSetOnclick(markers, onClick) {
        var _markers = [],//����ӵ���Ȥ���б�
          defer = $q.defer(),
          point,//��ǰ��ӵ������
          _length,//���鳤��
          _progress,//��ǰ������ӵĵ������
          _onClick;//����¼�����
        if (onClick) {
          _onClick = onClick;
        } else {
          _onClick = markerClick;
        }
        $timeout(function () {
          //�ж��Ƿ��ж�λ��
          if (!markers) {
            defer.reject('û�д�����Ȥ��');
            return;
          }
          //�����˲���
          if (!angular.isArray(markers)) {
            //����Ĳ���array
            if (markers.loc) {
              _markers.push(markers);
            } else {
              defer.reject('��ȡ����loc������Ϣ');
            }
          } else {
            if (markers[0].loc) {
              _markers = markers;
            } else {
              defer.reject('��ȡ����loc������Ϣ');
            }
          }
          _length = _markers.length - 1;
          angular.forEach(_markers, function (obj, index) {
            _progress = index;
            if (angular.isObject(obj.loc)) {
              point = new BMap.Point(obj.loc.lng, obj.loc.lat);
            } else if (angular.isString(obj.loc)) {
              point = new BMap.Point(obj.loc.split(',')[0], obj.loc.split(',')[1]);
            } else {
              _progress = '��' + index + '����Ȥ��loc���󲻴��ڻ��ʽ����ֻ֧��object��string';
            }
            var marker = new BMap.Marker(point);
            marker.obj = obj;
            marker.addEventListener('click', _onClick);
            getMap().addOverlay(marker);
            defer.notify(_progress);
            if (index === _length) {
              defer.resolve();
            }
          });
        });
        return defer.promise;
      }
    }
  }

  /***
   * ��ͼָ��
   */
  mapDirective.$inject = ['angularBMap'];
  function mapDirective(angularBMap) {
    return {
      restrict: 'EAC',
      replace: true,
      scope: true,
      template: '<div id="bMap" style="height: 100%;"></div>',
      link: mapLink,
      controller: mapController
    };
    /**
     * link
     * @constructor
     * @param scope
     * @param element
     * @param attr
     * @param ctrl
     */
    function mapLink(scope, element, attr, ctrl) {
      ctrl.initMap();
      ctrl.geoLocationAndCenter().then(function (result) {
        //��λ�ɹ�
        console.log(result);
      }, function (err) {
        //��λʧ��
        console.info(err);
      });
      var markers = [
        {loc: {lng: 121.496011, lat: 31.244085}},
        {lod: '121.494215,31.243005'},
        {loc: '121.493065,31.244981'},
        {lod: '121.49691,31.239454'},
        {loc: '121.502515,31.243622'}];
      ctrl.drawMarkersAndSetOnclick(markers).then(function (result) {
        console.log('��Ȥ��������');
      }, function (err) {
        console.log(err);
      }, function (progress) {
        console.log(progress);
      });
    }

    /**
     * controller
     * @constructor
     * @type {string[]}
     */
    function mapController() {
      this.geoLocation = angularBMap.geoLocation;//��λ
      this.initMap = angularBMap.initMap;//��ʼ��
      this.geoLocationAndCenter = angularBMap.geoLocationAndCenter;//��ȡ��ǰ��λ���ƶ�����ͼ����
      this.drawMarkers = angularBMap.drawMarkers;//�����Ȥ��
      this.drawMarkersAndSetOnclick = angularBMap.drawMarkersAndSetOnclick;//�����Ȥ��ͬʱ��ӵ���¼�
    }
  }
})(window, window.angular);