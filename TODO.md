### bug
* dependencyInjection scope prototype and request problem.(beego.orm,mysqlDao.Dao)

### feature
* ajkapi use function parameter name from kmgReflect
    * easy to test api (less code).
    * easy to write api (less code).
    * easy to call a api from another api (less code).
* put all session stuff in dependencyInjection.(no need to init session for enterScope Request,except )
* finish ProcessManager, apply it to dependencyInjection.
