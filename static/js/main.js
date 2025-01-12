

let apiUrl = 'http://localhost:3000/api/v1/gostore/store';

let app = document.getElementById("objects");

let fileLoading = document.getElementById("file-loading");

let backBtn = document.getElementById("back") 




async function getObjects(reesource, timeout = 5000, options = {}) {
    const response = await fetch(reesource, {
        ...options,
        signal: AbortSignal.timeout(timeout)
    })

    if (!response.ok) {
        const message = `An error has occured: ${response.status}`;
        throw new Error(message)
    }
    
    const objects = await response.json();

    return objects
}


function generateObjects(data) {

    let objects = []

    let dir = `
        <svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.5 8H4m0-2v13a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1V9a1 1 0 0 0-1-1h-5.032a1 1 0 0 1-.768-.36l-1.9-2.28a1 1 0 0 0-.768-.36H5a1 1 0 0 0-1 1Z"/>
        </svg> 
    `
    
    let file = `
        <svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
             <path stroke="currentColor" stroke-linejoin="round" stroke-width="2" d="M10 3v4a1 1 0 0 1-1 1H5m14-4v16a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V7.914a1 1 0 0 1 .293-.707l3.914-3.914A1 1 0 0 1 9.914 3H18a1 1 0 0 1 1 1Z"/>
        </svg>

    `

    data.forEach(element => {
        objects.push(
            `<li id="${element.name}" class="flex items-center space-x-3 rtl:space-x-reverse  hover:bg-gray-100 cursor-pointer">
                ${element.file ? file : dir}   
                <span>${element.name}</span>
            </li>
        `)
    });

    return objects
}

function removeBackground(data) {
    data.forEach(element => {
        let object = document.getElementById(element.name)
        
        object.classList.remove("bg-gray-100")
    })
}

function addBackEvent(element) {
    element.addEventListener("click", (event) => {
        console.log("hi")
    })
}

function addObjectEvents(data) {
    // For each element we add an event listener 
    console.log("I got called!")
    data.forEach(element => {
        let object = document.getElementById(element.name)
        
        object.addEventListener("dblclick", (event) => {
            console.log("Element " + element.name + " double clicked")
            
            let target = `/${element.name}`
            apiUrl = `${apiUrl}${target}`
            app.innerHTML = ""
            console.log(apiUrl)

            let objects = getObjects(apiUrl)
                .then(data => {
                    fileLoading.style = "display: none"
                    generateObjects(data).forEach(element => {
                        app.innerHTML += element
                    })
                    // Recursive
                    addObjectEvents(data)
                })
                .catch(error => {
                    console.log(error.message)
                });
            
            
        })

        object.addEventListener("click", (event) => {
            removeBackground(data)
            console.log("Element " + element.name + " Clicked once")
            object.classList.add("bg-gray-100")
        })

        
    })
}

objects = getObjects(apiUrl)
    .then(data => {
        fileLoading.style = "display: none"
        generateObjects(data).forEach(element => {
            app.innerHTML += element
        })

        addObjectEvents(data)
    })
    .catch(error => {
        console.log(error.message)
    });


// backBtn.addEventListener("click", (event) => {
//     let tmpApiUrl = apiUrl.replace(`/${apiUrl.split("/").slice(-1)}`, "")
//     console.log(tmpApiUrl)
//     if (tmpApiUrl.length > apiUrl) {
//         apiUrl = tmpApiUrl
//     }
    
//     objects = getObjects(apiUrl)
//     .then(data => {
//         fileLoading.style = "display: none"
//         generateObjects(data).forEach(element => {
//             app.innerHTML += element
//         })

//         addObjectEvents(data)
//     })
//     .catch(error => {
//         console.log(error.message)
//     });
    
// })