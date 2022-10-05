# blog
blog powered by golang gin web-framework and Gorm.
# startup by docker-compose
cd build/docker-compose &&  docker-compose up --build
![image](https://user-images.githubusercontent.com/9148977/194029154-867623e2-0c51-4075-b31e-723178c000a5.png)
![image](https://user-images.githubusercontent.com/9148977/194028300-6284e879-75c5-400e-a817-1fbb40164738.png)
![image](https://user-images.githubusercontent.com/9148977/194029317-a120a977-44fa-422a-aad6-90b1ec1c0637.png)



# how to build
docker build -t blog:v1 -f build/docker/Dockerfile .
![image](https://user-images.githubusercontent.com/9148977/194007818-f3392887-3aad-41f6-a9ba-734cc435965c.png)

# how to run the image
docker run -d  -p 8888:8888 blog:v1
![image](https://user-images.githubusercontent.com/9148977/194029423-96f36f39-ecb3-421c-832e-cf19ae8d1281.png)

