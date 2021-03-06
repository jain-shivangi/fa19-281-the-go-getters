import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpHeaders, HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-reviews',
  templateUrl: './reviews.component.html',
  styleUrls: ['./reviews.component.css']
})
export class ReviewsComponent implements OnInit {

  constructor(private http : HttpClient, private router: Router) { }

  reviews:any=[]
  ngOnInit() {
    console.log("*****")
    if(sessionStorage.getItem('userId')==null)
    {
      this.router.navigate(['./login'])
      window.alert("you need to login first!")
    }
    this.router.navigate(['./reviews'])
    this.username=sessionStorage.getItem('username')
  }
 endpoint:any="https://i18253eej8.execute-api.us-east-1.amazonaws.com/prod"
 itemName:any
 reviewForAnItem:any={}
 ItemNameEntered:any
 CommentsEntered:any
 RatingEntered:any
 RequestObject={}
 ReviewPost=[]
 username:any
 bool1:any=false
  getReview(){

    let header = new HttpHeaders();
    header.append('Content-Type', 'application/json');
    this.bool1=true
     this.http
        .get(this.endpoint+'/getReviews'+'/'+this.itemName,{headers: header})
        .subscribe((res) => {
          this.reviews=res
          console.log(this.reviews)
            //do something with the response here
            this.router.navigate(['./reviews']);
            console.log(res);

});
}
postReview(){

  this.createObject()
  let header = new HttpHeaders();
  header.append('Content-Type', 'application/json');
  
   this.http
      .post(this.endpoint+'/postReview',JSON.stringify(this.RequestObject),{headers: header})
      .subscribe((res) => {
          //do something with the response here
          this.router.navigate(['./reviews']);
          console.log(res);

});
}
createObject(){
this.ReviewPost[0]={'ReviewerName':this.username,"Comment":this.CommentsEntered,"Rating":parseFloat(this.RatingEntered)}
console.log(this.ReviewPost)
this.RequestObject={"ItemName":this.ItemNameEntered,"Reviews":this.ReviewPost}
console.log(this.RequestObject)
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
