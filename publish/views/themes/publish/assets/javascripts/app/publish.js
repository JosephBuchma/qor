(function (factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as anonymous module.
    define(['jquery'], factory);
  } else if (typeof exports === 'object') {
    // Node / CommonJS
    factory(require('jquery'));
  } else {
    // Browser globals.
    factory(jQuery);
  }
})(function ($) {

  'use strict';

  var NAMESPACE = 'qor.publish';
  var EVENT_CLICK = 'click.' + NAMESPACE;
  var EVENT_SHOWN = 'shown.bs.modal';
  var EVENT_HIDDEN = 'hidden.bs.modal';

  function Publish(element, options) {
    this.$element = $(element);
    this.options = $.extend({}, Publish.DEFAULTS, $.isPlainObject(options) && options);
    this.loading = false;
    this.init();
  }

  Publish.prototype = {
    constructor: Publish,

    init: function () {
      var options = this.options;
      var $this = this.$element;

      this.$modal = $(Publish.MODAL).appendTo('body');
      this.bind();
    },

    bind: function () {
      this.$element.on(EVENT_CLICK, $.proxy(this.click, this));
      this.$modal.
        on(EVENT_SHOWN, function () {
          $(this).trigger('enable.qor.textviewer');
        }).
        on(EVENT_HIDDEN, function () {
          $(this).trigger('disable.qor.textviewer');
        });
    },

    unbind: function () {
      this.$element.off(EVENT_CLICK, this.click);
      this.$modal.off(EVENT_SHOWN).off(EVENT_HIDDEN);
    },

    click: function (e) {
      var options = this.options;
      var $target = $(e.target);

      if ($target.is(options.toggleView)) {
        e.preventDefault();

        if (this.loading) {
          return;
        }

        this.loading = true;
        this.$modal.find('.mdl-card__supporting-text').empty().load($target.data('url'), $.proxy(this.show, this));
      } else if ($target.is(options.toggleCheck)) {
        if (!$target.prop('disabled')) {
          $target.closest('table').find('tbody :checkbox').click();
        }
      }
    },

    show: function () {
      this.loading = false;
      this.$modal.qorModal('show');
    },

    destroy: function () {
      this.unbind();
      this.$element.removeData(NAMESPACE);
    },
  };

  Publish.DEFAULTS = {
    toggleView: '.qor-action__view',
    toggleCheck: '.qor-action__check-all',
  };

  Publish.MODAL = (
    '<div class="qor-modal fade" tabindex="-1" role="dialog" aria-hidden="true">' +
      '<div class="mdl-card mdl-shadow--2dp" role="document">' +
        '<div class="mdl-card__title mdl-card--border">' +
          '<h2 class="mdl-card__title-text">Changes</h2>' +
        '</div>' +
        '<div class="mdl-card__supporting-text"></div>' +
        '<div class="mdl-card__actions mdl-card--border">' +
          '<a class="mdl-button mdl-button--colored mdl-js-button mdl-js-ripple-effect" data-dismiss="modal">Close</a>' +
        '</div>' +
        '<div class="mdl-card__menu">' +
          '<button class="mdl-button mdl-button--icon mdl-js-button mdl-js-ripple-effect" data-dismiss="modal" aria-label="close">' +
            '<i class="material-icons">close</i>' +
          '</button>' +
        '</div>' +
      '</div>' +
    '</div>'
  );

  Publish.plugin = function (options) {
    return this.each(function () {
      var $this = $(this);
      var data = $this.data(NAMESPACE);
      var fn;

      if (!data) {
        $this.data(NAMESPACE, (data = new Publish(this, options)));
      }

      if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
        fn.apply(data);
      }
    });
  };

  $(function () {
    Publish.plugin.call($('.qor-table'));
  });

});
