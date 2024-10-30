use studentgrpc::student_client::StudentClient;
use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use studentgrpc::StudentRequest;
use serde::{Deserialize, Serialize};

pub mod estudiante {
    tonic::include_proto!("estudiante")
}

#[derive(Deserialize,Serialize)]
struct StudentData {
    name: String,
    age: i32,
    faculty: String,
    discipline: i32,
}


async fn handle_student(student: web::Json<StudentData>) -> impl Responder {
    let mut client = match StudentClient::connect("http://go-server-service:50051").await {
        Ok(client) => client,
        Err(e) => return HttpResponse::InternalServerError().body(format!("Failed to connect to gRPC server: {}", e)),
    };

    let request = tonic::Request::new(StudentRequest {
        name: student.name.clone(),
        age: student.age,
        faculty: student.faculty.clone(),
        discipline: student.discipline,
    });

    match client.get_student(request).await {
        Ok(response) => {

            println!("RESPONSE={:?}", response);

            HttpResponse::Ok().json
            (format!("Student: {:?}", response))
        },
        Err(e) => HttpResponse::InternalServerError().body(format!("gRPC call failed: {}", e)),
    }

}



#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Starting server at http://localhost:8080");
    HttpServer::new(|| {
        App::new()
            .route("/faculty", web::post().to(handle_student))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}


fn main() {
    println!("Hello, world!");
}
