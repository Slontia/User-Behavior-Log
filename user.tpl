<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">

    <title>Bootstrap Template Page for Go Web Programming</title>

    <link href="//cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
  </head>

  <body>
    <p><a class="btn btn-primary btn-lg" href="#" role="button" id="bbb">Send &raquo;</a></p>

    <script src="//cdn.bootcss.com/jquery/1.10.2/jquery.min.js"></script>  
    <script src="//cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
  </body>
  <script>
    bindEvent()
    function bindEvent()
    {
      $("#bbb").click(function(){
        var data_dict = {
          'wahr': "7777",
          'psw': "fuU"
        }
        var ajax_data = {
          'timestp': "123",
          'sign': "aaa",
          'action': data_dict
        }
        $.ajax({
          //ContentType: 'application/json; charset=utf-8',
          //dataType: 'text',
          url: '/entry/',
          type: 'POST',
          data: JSON.stringify(ajax_data),
          //traditional: true,
          success: function(d){
            if (d === "")
            {
              alert("Success");
            }
            else
            {
              alert(d)
            }
          },
          error: function(){
            alert("Lost Connection")
          }
        });
      })
    }
  </script>
</html>