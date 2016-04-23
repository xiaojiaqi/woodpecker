/* global Lib, Timesheet */

(function(){
  'use strict';
  
  Lib.ready(function() {
    console.log('ads');
    
    /* jshint -W031 */
    new Timesheet('timesheet-default', 1900, 1927, [
        ['1903', '7/1905', 'xxxx', 'lorem']   ,
      ['6/1905', '09/1925', 'Some great memories', 'ipsum'],
      ['1910', 'Had very bad luck'],
      ['10/1923', '1916', 'At least had fun', 'dolor'],
      ['2/1925', '05/1926', 'Enjoyed those times as well', 'ipsum'],
      ['7/1920', '09/1927', 'Bad luck again', 'default'],
      ['10/1918', '1927', 'For a long time nothing happened', 'dolor'],
      ['01/1908', '05/1927', 'LOST Season #4', 'lorem'],
      ['01/1915', '05/1927', 'LOST Season #4', 'sit'],
      ['02/1901', '05/1920', 'LOST Season #5', 'lorem'],
      ['09/1908', '06/1927', 'FRINGE #1 & #2', 'ipsum']
    ]);
    
 

    document.querySelector('#switch-dark').addEventListener('click', function() {
      document.querySelector('body').className = 'index black';
    });

    document.querySelector('#switch-light').addEventListener('click', function() {
      document.querySelector('body').className = 'index white';
    });
  });
})();
