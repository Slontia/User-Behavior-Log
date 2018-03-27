<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">

    <title>Bootstrap Template Page for Go Web Programming</title>

    <link href="//cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
  </head>

  <body>
    <p>全部数据：</p>
    <p id="datatxt"></p>  

    <script src="//cdn.bootcss.com/jquery/1.10.2/jquery.min.js"></script>  
    <script src="//cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
  </body>
  <script>
    // function readData()
    // {
    //   $.ajax({
    //     url: '/show/',
    //     type: 'GET',
    //     success: function(d) {
    //       document.getElementById("datatxt").innerText = d;
    //     },
    //     error: function(d) {
    //       alert("Lost Connection");
    //     }
    //   })
    // }

    function bindEvent()
    {
      $("#sendbtn").click(function(){
        var ajax_data = {
          'time': "123",
          'sign': "aaa",
          'action': "login",
          'id': "xiaoming"
        }
        $.ajax({
          url: '/entry/',
          type: 'POST',
          data: JSON.stringify(ajax_data),
          success: function(d){
            if (d === "")
            {
              alert("Success");
            }
            else
            {
              alert(d);
            }
          },
          error: function(){
            alert("Lost Connection");
          }
        });
      })
    }
  </script>
</html>