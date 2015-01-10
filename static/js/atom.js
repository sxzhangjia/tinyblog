/*

project: atomjs

*/

(function(w){

	var d = w.document,
		l = w.location;

	var config = {
		basePath : "",
		jsFileTail : ".js",
		cssFileTail : ".css",
		modAttrNameKey : "data-name",
		modAliasNameKey : "data-alias"
	};

	//
	var runTime = {};

	//storge loaded module object
	var moduleMap = {};

	//
	var loadingMap = {};

	//
	var moduleCache = null;

	//
	var Utils = {
		getJSIntactURL: function (modName) {
			return config.basePath + modName + config.jsFileTail;
		},
		getSelfElem: function () {
			var s = d.getElementsByTagName('script');
			return s[s.length-1];
		},
		getBasePath: function (s) {
			return s.replace(/[^\/]+\.js/, '');
		}
	};

	var Fn = {
		getModuleCache: function () {
			var c = moduleCache;
			moduleCache = null;
			return c;
		},
		setModuleCache: function (v) {
			moduleCache = v;
		},
		getReadyModule: function (k) {
			return moduleMap[k];
		},
		setReadyModule: function (k, v) {
			//console.log(k, v);
			moduleMap[k] = v;
		}
	};

	function addScript(url, attr, loaded) {

		var s = d.createElement("script");
		s.type = "text/javascript";

		for(var k in attr) {
			s.setAttribute(k, attr[k]);
		}
		
		s.onload = loaded;
		s.src = url;
		d.body.appendChild(s);
	}

	function moduleLoaded() {

		var name = this.getAttribute(config.modAttrNameKey);
		var alias = this.getAttribute(config.modAliasNameKey);
		var cache = Fn.getModuleCache();
		var queue = loadingMap[name];

		cache = cache.f.call(null, cache.g, cache.m);
		Fn.setReadyModule(name, cache);

		for(var loader = queue.shift(); loader; ) {
			loader.loaded(name, alias, cache);
			loader = queue.shift();
		}
	}

	function loadModule(name, alias, loader) {

		var loadedMod = Fn.getReadyModule(name);
		if(loader && loadedMod) {
			loader.loaded(name, alias, loadedMod);
			return;
		}

		var loadingMod = loadingMap[name];
		if(loader && loadingMod) {
			loadingMod.push(loader);
		} else {
			loadingMap[name] = [loader];
		}

		var attr = {};
		attr[config.modAttrNameKey] = name;
		attr[config.modAliasNameKey] = alias;
		addScript(Utils.getJSIntactURL(name), attr, moduleLoaded);
	}

	function DepsLoader(deps, factory) {
		
		var self = this;

		self.depsMap = {};
		self.num = 0;
		self.loadedNum = 0;

		self.factory = factory;
		self.deps = deps;

		self.onload = function () {};
	}

	DepsLoader.prototype = {

		load: function () {

			var self = this;
			var deps = self.deps;
			for(var alias in deps) {
				self.num++;
				loadModule(deps[alias], alias, self);
			}
		},

		loaded: function (modName, alias, mod) {

			var self = this;
			self.depsMap[alias] = mod;
			self.loadedNum++;
			if(self.num === self.loadedNum) {
				self.onload(self.depsMap);
			}
		}

	};


	w.define = function(deps, factory) {
		
		var loader = null;

		if(typeof deps === 'function') {
			
			factory = deps;

		} else {

			loader = new DepsLoader(deps, factory);
			loader.onload = function(m) {

				//console.log(m, new Date().getTime())
				//Fn.setModuleCache(this.factory(runTime, m));

			};
			loader.load();

		}

		Fn.setModuleCache({
			f: factory,
			g: runTime,
			m: {}
		});

	};

	w.require = function(deps, factory) {

		

	};

	//init
	function init() {

		var self = Utils.getSelfElem(),
			selfUrl = self.src,
			mainModName = self.getAttribute('data-main');

		self = null;

		config.basePath = Utils.getBasePath(selfUrl);

		//addScript(Utils.getJSIntactURL(mainModName), {});
		loadModule(mainModName, 'main');
	}

	init();

	//debug
	w.moduleMap = moduleMap;
	w.loadingMap = loadingMap;

})(window);