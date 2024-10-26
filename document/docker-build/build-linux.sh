#!/bin/bash

REPO="DAIRO-HY/DairoNPS"
BRANCH="release"

#删除上次编译文件
rm DairoNPS.zip
rm -rf DairoNPS-main
rm /app/DairoNPS/dairo-nps-linux-amd64

curl -L -o DairoNPS.zip https://github.com/DAIRO-HY/DairoNPS/archive/refs/heads/main.zip
unzip DairoNPS.zip
cd DairoNPS-main

#开始编译 由于使用了sqlite插件，编译时需要指定参数CGO_ENABLED=1
CGO_ENABLED=1 go build -o /app/DairoNPS/dairo-nps-linux-amd64

/app/DairoNPS/dairo-nps-linux-amd64







pipeline {
    agent any
    tools{
        go "go-1.23.1"
    }
    environment {
        REPO = "DAIRO-HY/DairoNPS"
        BRANCH = "release"
    }
    stages {
        stage("Test") {
            steps {
                sh 'go version'
            }
        }
        stage("拉取代码") {
            steps {
                script{
                    if(!fileExists("DairoNPS")){
                        def cloneUrl = "https://${env.GITHUB_TOKEN}@github.com/${env.REPO}.git"
                        echo "克隆地址:${cloneUrl}"
                        sh "git clone --branch ${env.BRANCH} ${cloneUrl}"
                    }
                }
                dir("DairoNPS"){
                    sh "git pull"
                }
            }
        }
        stage("编译打包") {
            steps {
                dir("DairoNPS"){
                    sh "CGO_ENABLED=1 go build -o dairo-nps-linux-amd64"
                    sh "ls"
                }
            }
        }
        /**stage("获取版本号") {
            steps {
                dir("DairoDfs/dairo-dfs-server"){
                    script{
                        def buildGradleTxt = sh(script: "cat build.gradle.kts", returnStdout: true)
                        def pattern = ~/version = \"(.*)\"/
                        def matcher = pattern.matcher(buildGradleTxt)
                        matcher.find()
                        def version = matcher.group(1)
                        env.VERSION = version
                        env.TAG_NAME = version
                    }
                }
            }
        }
        stage("推送代码") {
            steps {
                dir("DairoDfs"){
                    script{
                        try{
                            sh "git tag -d ${env.TAG_NAME}"//先删除已经存在的tag,如果存在
                        }catch(Exception e){
                            //e.printStackTrace()
                        }
                    }
                    script{
                        try{
                            sh "git push origin --delete tag ${env.TAG_NAME}"//先删除已经存在的tag,如果存在
                        }catch(Exception e){
                            //e.printStackTrace()
                        }
                    }

                    sh "git tag ${env.TAG_NAME}"
                    sh "git push origin ${env.TAG_NAME}"
                }
            }
        }

        stage("创建标签") {
            steps {
                script {
                    def releaseBody = "本次发布版本${env.TAG_NAME}"
                    echo releaseBody

                    def releaseScript = """
                        curl -L -X POST "https://api.github.com/repos/${env.REPO}/releases" \
                        -H "Accept: application/vnd.github.v3+json" \
                        -H "Authorization: Bearer ${env.GITHUB_TOKEN}" \
                        -H "X-GitHub-Api-Version: 2022-11-28" \
                        -d "{\\\"tag_name\\\":\\\"${env.TAG_NAME}\\\",\\\"name\\\":\\\"${env.TAG_NAME}\\\",\\\"body\\\":\\\"${releaseBody}\\\"}"
                        """
                    def response = sh(script: releaseScript, returnStdout: true).trim()
                    echo "-->result:${response}"
                    def releaseId = readJSON(text: response).id
                    env.RELEASE_ID = releaseId
                }
            }
        }

        stage("上传jar包") {
            steps {
                dir("DairoDfs/dairo-dfs-server/build/libs/"){
                    script {
                        def filePath = "dairo-dfs-server-${env.VERSION}.jar"
                        //def fileName = sh(script: "basename ${filePath}", returnStdout: true).trim()
                        def response = sh(script: """
                            curl -s -X POST \
                            -H "Accept: application/vnd.github+json" \
                            -H "Authorization: Bearer ${env.GITHUB_TOKEN}" \
                            -H "X-GitHub-Api-Version: 2022-11-28" \
                            -H "Content-Type: application/octet-stream" \
                            --data-binary "@${filePath}" \
                            "https://uploads.github.com/repos/${env.REPO}/releases/${env.RELEASE_ID}/assets?name=dairo-dfs-server.jar" \
                            """, returnStdout: true).trim()
                        echo "-->result:${response}"
                    }
                }
            }
        }*/
        /**
        stage("上传docker镜像") {
            steps {
                sh "cp DairoDfs/document/docker/Dockerfile DairoDfs/dairo-dfs-server/build/libs/"
                dir("DairoDfs/dairo-dfs-server/build/libs/"){

                    //重命名
                    sh "mv dairo-dfs-server-${env.VERSION}.jar dairo-dfs-server.jar"
                    sh "ls"
                    sh "docker --version"
                    sh "docker build -t ${env.DOCKER_USER}/dairo-dfs:${env.VERSION} ."
                    sh "docker login -u ${env.DOCKER_USER} --password ${env.DOCKER_PASSWORD}"
                    sh "docker push dairopapa/dairo-dfs:${env.VERSION}"
                    sh "docker logout"
                }
            }
        }
        */
    }
}
