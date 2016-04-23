
		(function(){
	  'use strict';

	  Lib.ready(function() {
	    console.log('ads');

	   // replace code
	   new Timesheet('timesheet-default', 1900, 1927, 
		   [['1/1900', '1/1900', '1-E:/fakewechat/src/github.com/woodpecker/client2.go:23 main.oneclient', 'lorem'],
['1/1900', '1/1913', '2.1-E:/fakewechat/src/github.com/woodpecker/client2.go:72 main.A', 'ipsum'],
['1/1913', '1/1913', '2.2-E:/fakewechat/src/github.com/woodpecker/client2.go:75 main.A', 'dolor'],
['1/1913', '1/1926', '2.2-E:/fakewechat/src/github.com/woodpecker/client2.go:77 main.A', 'ipsum'],
['1/1913', '1/1913', '2.3.1-E:/fakewechat/src/github.com/woodpecker/client2.go:86 main.B', 'default'],
['1/1913', '1/1913', '2.3.2-E:/fakewechat/src/github.com/woodpecker/client2.go:90 main.B', 'sit'],
['1/1926', '1/1926', '2.3.3.1-E:/fakewechat/src/github.com/woodpecker/client2.go:101 main.C', 'lorem'],
['1/1926', '1/1926', '2.3.3.1-E:/fakewechat/src/github.com/woodpecker/client2.go:105 main.C', 'ipsum']
	 ]);
	 



	    document.querySelector('#switch-dark').addEventListener('click', function() {
	      document.querySelector('body').className = 'index black';
	    });

	    document.querySelector('#switch-light').addEventListener('click', function() {
	      document.querySelector('body').className = 'index white';
	    });
	  });
	})();
    