import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpHeaders, HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-registration',
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.css']
})
export class RegistrationComponent implements OnInit {
  
  username:any ='';
  password:any='';
  firstname:any='';
  lastname:any='';
  userId:any


  constructor(private http : HttpClient, private router: Router) { }
  ngOnInit() {
  }

  signup(){

    let header = new HttpHeaders();
    header.append('Content-Type', 'application/json');
    
     this.http
        .post("https://i18253eej8.execute-api.us-east-1.amazonaws.com/prod/signup",
        JSON.stringify({"username":this.username,"password":this.password,"firstname":this.firstname,"lastname":this.lastname}),
            {headers: header})
        .subscribe((res) => {

            //do something with the response here
            if(res['id']==null)
            {
              window.alert("User already exists")
              this.router.navigate(['./signup'])
              return
        }
            sessionStorage.setItem('userId',res['id']);
            this.userId=sessionStorage.getItem('userId')
            this.router.navigate(['./home']);


            console.log(res);

});
}
gotoMenu(){
  this.router.navigate(['./menu'])
}

gotoReviews(){
  this.router.navigate(['./reviews'])
}

gotoHome(){
  this.router.navigate(['./home'])
}
logout(){
  sessionStorage.setItem('userId',null)
  this.router.navigate(['./login'])
}
}
