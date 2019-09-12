def setBuildStatus(String message, String state, String project_url, String sha) {
    step([
        $class: "GitHubCommitStatusSetter",
        reposSource: [$class: "ManuallyEnteredRepositorySource", url: project_url],
        contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/build-status"],
        errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
        commitShaSource: [$class: "ManuallyEnteredShaSource", sha: sha ],
        statusBackrefSource: [$class: "ManuallyEnteredBackrefSource", backref: project_url + "/commit/" + sha],
        statusResultSource: [$class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
    ]);
}

def archiveAndCleanup() {
    dir (env.WORKSPACE) {
        junit allowEmptyResults: true, testResults: "pkg/**/test-report.xml,test/**/test-report.xml,build/_output/test-logs/operator-e2e-local-test.xml,build/_output/test-logs/operator-e2e-helm-test.xml,build/_output/test-logs/operator-e2e-test.xml,java/**/surefire-reports/*.xml,java/**/failsafe-reports/*.xml"
        archiveArtifacts onlyIfSuccessful: false, allowEmptyArchive: true, artifacts: 'build/**/*,deploy/**/*,java/utils/target/test-output/**/*,java/utils/target/surefire-reports/**/*,java/utils/target/failsafe-reports/**/*,java/functional-tests/target/test-output/**/*,java/functional-tests/target/surefire-reports/**/*,java/functional-tests/target/failsafe-reports/**/*'
        sh '''
            helm delete --purge $(helm ls --namespace $TEST_NAMESPACE --short) || true
            kubectl delete clusterrole $TEST_NAMESPACE-coherence-operator || true
            kubectl delete clusterrolebinding $TEST_NAMESPACE-coherence-operator-cluster || true
            kubectl delete namespace $TEST_NAMESPACE --force --grace-period=0 || true
            make uninstall-crds || true
        '''
    }
}

pipeline {
    agent {
        label 'Kubernetes'
    }
    environment {
        HTTP_PROXY  = credentials('coherence-operator-http-proxy')
        HTTPS_PROXY = credentials('coherence-operator-https-proxy')
        NO_PROXY    = credentials('coherence-operator-no-proxy')
        PROJECT_URL = "https://github.com/oracle/coherence-operator"

        COHERENCE_IMAGE_PREFIX = credentials('coherence-operator-coherence-image-prefix')
        TEST_IMAGE_PREFIX      = credentials('coherence-operator-test-image-prefix')

        TEST_NAMESPACE = "test-cop-${env.BUILD_NUMBER}"
    }
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '28', numToKeepStr: '')
        timeout(time: 4, unit: 'HOURS')
    }
    stages {
        stage('code-review') {
            steps {
                echo 'Code Review'
                script {
                    setBuildStatus("Code Review in Progress...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                '''
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                    make code-review
                    '''
                }
            }
        }
        stage('build') {
            steps {
                echo 'Build'
                script {
                    setBuildStatus("Build in Progress...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                '''
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                    make clean
                    export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                    export TEST_MANIFEST_VALUES=deploy/oci-values.yaml
                    make build-all
                    '''
                }
            }
        }
        stage('release') {
            when {
                expression { env.RELEASE_ON_SUCCESS == 'true' }
            }
            steps {
                echo 'Release'
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                '''
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                    export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                    git config user.name "Coherence Bot"
                    git config user.email coherence-bot_ww@oracle.com
                    make clean
                    make build-all-images
                    make release RELEASE_DRY_RUN=${DRY_RUN} RELEASE_IMAGE_PREFIX=${RELEASE_IMAGE_REPO} VERSION_SUFFIX=${RELEASE_SUFFIX}
                    '''
                }
            }
        }
    }
    post {
        always {
            script {
                archiveAndCleanup()
            }
            deleteDir()
        }
        success {
            setBuildStatus("Build succeeded", "SUCCESS", "${env.PROJECT_URL}", "${env.GIT_COMMIT}");
        }
        failure {
            setBuildStatus("Build failed", "FAILURE", "${env.PROJECT_URL}", "${env.GIT_COMMIT}");
        }
    }
}
