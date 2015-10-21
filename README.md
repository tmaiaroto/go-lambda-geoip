This is a simple example of Go for AWS Lambda using a Node.js wrapper. Follow along here and you should have 
a pretty fun IP -> geoloacation service at the end of it. A really inexpensive service at that. This is a 
follow up to [an article I wrote here.](https://medium.com/@shift8creative/go-amazon-lambda-7e95a147cec8#.ab93bgu8s)

## Instructions

### Prerequisites
First, you need an AWS account of course. After that, you'll need to download [Maxmind's GeoLite2 City database.](http://dev.maxmind.com/geoip/geoip2/geolite2/) 
The filename will likely be ```GeoLite2-City.mmdb```.

You can't run the Go program locally unfortunately. Though I'm going to see if there's a way to send the input 
required and simulate Lambda. For now, you'll just need to run it in Lambda.

I've included a pre-compiled binary for you, but you'd normally then compile for Linux.

```
GOOS=linux GOARCH=amd64 go build main.go
```

### Creating the AWS Lambda Function
I've included an ```index.js``` file as well that will call the binary. You'll want to zip those 3 files up 
and upload to a new AWS Lambda function. Again the files to zip are:

```
index.js    
main    
GeoLite2-City.mmdb
```

You'll then need to ensure the handler is ```index.handler``` (which is the default AWS Lambda uses - you can 
change this if you really want, it's simply the file name and exported method). 

A basic execution role will do. Under the Advanced section, you won't need to allocate anymore than 128MB of RAM 
and a timeout of a few second should be enough. I actually set 10, because there's no harm, but this really runs 
in milliseconds not seconds. You are billed for actual duration and it rounds up to the nearest 100ms.

### Exposing via API Gateway 
Great, now you'll just need a way to invoke it with the required parameters. We'll use Amazon API Gateway.

Go to the API endpoints tab and add a new endpoint. Call the API whatever you want, use whatever deployment 
stage you like, but the method should be GET. Mine ended up being ```/prod/test-go``` for example.

I'm not a huge fan of API Gateway's dashboard, but hang in there. Here's where it gets complex...

If you're just playing around, you're likely not going to use your own domain name or anything and you'll
just use the default AWS assigns you. Cool...But you'll also likely want to disable the need for an API key 
at that point too.

To do this, go to the "Dashboard" dropdown menu item and then "Resources" item. You should see your methods. 
The will likely be collapsed so expand the tree on the left. It starts with "/". If not, ensure you are looking at 
the right API if you have multiple, it's the dropdown menu next to Resources.

Go down to the last item in the tree for the API endpoint you created for this function. It'll be "GET" with a single
cube icon. Now you'll see a visual thing on the right that's supposed to illustrate a request to your API.

Click "Method Request" and you'll see the "Authorization Settings" there. Make sure you don't need an API key (if this
is what you want).

Cool. Now the most important part. ***Mapping templates.***

Under the "Method Execution" (where you were previously before setting the authorization stuff) you'll see on the right 
"Integration Request". Click that. It's going to show you the Lambda function the endpoint invokes and at the bottom 
you'll see "Mapping Templates".

Under here add a new one for ```application/json``` Content-Type (which is the example helper text in the input).
Click the little checkmark icon. Then click the text in the list on the left that you just created. Now to the right 
will be a confusing empty area. Click the pencil icon next to "Input passthrough" and it'll turn into a select dropdown. 
You'll want to change it to "Mapping template" and then you'll see a little code editor area with numbered lines.

You can create models and re-use them. So in the future you can save a small step. Perhaps complex step based on what
you're doing. What you'll want to do is simply paste the following in there:

```
{
"stage" : "$context.stage",
"request-id" : "$context.requestId",
"api-id" : "$context.apiId",
"resource-path" : "$context.resourcePath",
"resource-id" : "$context.resourceId",
"http-method" : "$context.httpMethod",
"source-ip" : "$context.identity.sourceIp",
"user-agent" : "$context.identity.userAgent",
"account-id" : "$context.identity.accountId",
"api-key" : "$context.identity.apiKey",
"caller" : "$context.identity.caller",
"user" : "$context.identity.user",
"user-arn" : "$context.identity.userArn"
}
```

Then click the little checkmark icon inside a circle up top next to that "Mapping template" select option you chose.

You technically don't need all of these. You really only need ```source-ip``` for this example. But so you know what's 
available to you - there you are. Yes, this is documented somewhere under AWS documentation but what I just walked you 
through is probably much easier than going through the documentation. At least it took me a while to figure this out, 
but it's really quite simple. It's just annoying to configure all the time.

Anyway, now your Lambda function has the data it needs. Click "Deploy API" up top the left there.

### Try Out Your API
Go back up top to the menu item which now says "Resources" and switch it back to the "Dashboard" item. You'll be taken
to the overview with some graphs for your calls and such. Here in a blue info box you'll see your API URL.

Add to that the endpoint path because it's just the domain and stage path. 

You should see a JSON response with your gelocation based on your IP address. Note that it may not be super accurate
because you're using Maxmind's free data set. You can pay them some more money to get better accuracy...But at this 
point you should see how to use Go with AWS Lambda (ignoring all the hocus-pocus).

## Hocus-Pocus How?
The method to spawn a Go process via Node.js can be found here: 
https://github.com/jasonmoo/lambda_proc

It's README contains a little explanation of what's going on. There are other ways to spawn a Go process with Node.js,
but this was a very clever way to handle it and it does make a huge difference for performance.

This is important because you are billed based on duration and this will execute faster (after the frist call from a 
"cold" start).

Until Amazon decides to support Go natively, we'll need to use Node.js or perhaps even Python or Java to execute our
Go binary. Fortunately this works. It's also nice to note here that any assets you may need, such as that database file
from Maxmind, can be zipped up and remain accessible to your Go app.

There's a fair amount of configuration involved to set this up. You could use the JAWS framework to help make some of 
that faster. It's an AWS Lambda framework that will deploy from your command line. However, you will need to back and
zip those files up and upload them yourself because it doesn't take into consideration your Go app. It's a more than 
just a tool, it is indeed a framework. So it's great hocus-pocus too, but still leaves you with some work.

## Performance Considerations, Observations, and Thoughts
I did get away uploading the zip file via the console despite it suggesting to use S3 for the larger file. 

I didn't seem to need more RAM than the size of the database file (makes sense). 

I would like to see what would happen if I stored the database in a memory mapped file (like BoltDB). This may 
help matters when working with larger data sets. Of course there's always DynamoDB or RDS or something...But what's
going to be faster and easier for deployment?

I threw some load at the API, about 5,000 requests 10 at a time. It averaged 46 some requests per second. Between 
all my testing and this light load test, I think I was billed $0.03. So Lambda is quite inexpensive.

The fastest response was 0.1677 seconds. The average was 0.2146 and of course your milage will vary. There's many
things to consider, including the region and your connection to the internet.

While the average response was 200 some milliseconds, keep in mind this includes API Gateway. I have seen this 
Lambda function execute in under one millisecond when using test data and invoking from the AWS Lambda dashboard.

Even though we're talking nanoseconds in some cases, it still doesn't matter because we're billed for 100ms minimum.
Well, _it doesn't matter that much_...It matters in the sense that perhaps this function in Go executed under 100ms, 
while one in Node.js executed anywhere over 100ms. Then over enough requests, it will maybe be noticably cheaper to use Go. 
Maybe. Depends on how many requests. I should be fair in mentioning that, though it won't matter for many.

The Lambda function itself can execute quite quickly once warm, but the API Gateway can take a bit longer to 
deliver your response. However, it's still quite fast. I do believe Lambda + API Gateway is a viable solution.