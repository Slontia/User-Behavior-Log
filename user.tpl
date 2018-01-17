<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>Bootstrap Template Page for Go Web Programming</title>

    <!-- Bootstrap core CSS -->
    <!-->
    <link href="bootstrap/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="jquery/dist/jquery-1.11.2.min.js"></script>
    <script src="bootstrap/dist/js/bootstrap.min.js"></script>
    <-->


    
    <link href="//cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
  </head>

  <body>

    <nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container">
        <div class="navbar-header">
          <a class="navbar-brand" href="#">Person general infor</a>
        </div>
      </div>
    </nav>

    <div class="jumbotron">
      <div class="container">
        <h1>Hello, {{.Name}}</h1>
        <ul>
        <li>Name   : {{.Name}}<p>
        <li>Id     : {{.Id}}<p>
        <li>Country: {{.Country}}
        </ul>
        <p><a class="btn btn-primary btn-lg" href="#" role="button" id="bbb">More &raquo;</a></p>
        <!--button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button-->
      </div>
    </div>
<!-->
    <div class="container">
      <div class="row">
        <div class="col-md-4">
          <h2>Name</h2>
          <p>Name has the value of : {{.Name}} </p>
          <p><a class="btn btn-default" href="#" role="button">More &raquo;</a></p>
        </div>
        <div class="col-md-4">
          <h2>Id</h2>
          <p>Id has the value of : {{.Id}} </p>
          <p><a class="btn btn-default" href="#" role="button">More &raquo;</a></p>
       </div>
        <div class="col-md-4">
          <h2>Country</h2>
          <p>Country has the value of : {{.Country}} </p>
          <p><a class="btn btn-default" href="#" role="button">More &raquo;</a></p>
        </div>
      </div>
<-->
      <hr>

      <footer>
      <nav class="navbar navbar-inverse ">
        <div class="container">
          <div class="navbar-header">
            <a class="navbar-brand" href="#">Hello, {{.Name}}, Welcome to go programming...</a>
          </div>
        </div>
      </nav>
      </footer>
    </div> <!-- /container -->

    <script src="//cdn.bootcss.com/jquery/1.10.2/jquery.min.js"></script>  
    <script src="//cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
  </body>
  <script>
    addshanchu()
    //alert(999)
  function addshanchu()
  {

    $("#bbb").click(function(){
    /*var xhr = new XMLHttpRequest();
    xhr.open('get', '/ajax', true);
    xhr.send();*/
    /*var data_dict = new Array()
    data_dict["name"] = "xiaoming"
    data_dict["psw"] = "fuQ"
    */
    var data_dict = {
      'name': "xiaoming",
      'psw': "fuQ"
    }
    var ajax_data = {
      'api': "Login",
      'timestp': "123",
      'sign': "aaa",
      'data': data_dict
    }
    alert(ajax_data)
      $.ajax({
        //ContentType: 'application/json; charset=utf-8',
        //dataType: 'text',
        url: '/entry/',
        type: 'POST',
        //data: ajax_data,
        data: JSON.stringify(ajax_data),
        //traditional: true,
        success: function(d){
          alert(d.Message)
          if(d.trim()=="OK")
          {
            alert("删除成功");
          }
          else
          {
            alert("删除失败");
          }
        },
        error: function(){
          alert("错误")
        }
      });
    })
  }
  </script>
</html>