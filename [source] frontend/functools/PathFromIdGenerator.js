export default function PathFromIdGenerator(id) {
    let result = ""
    try{
        id = id.toString()
        const idLen = id.length
        const twoLetterSize = idLen >> 1

        for (let i= 0; i < twoLetterSize; i++){
            result += `/${id[idLen - 2*i - 1]}/${id[idLen - 2*i - 2]}`
        }

        if (idLen !== twoLetterSize << 1){
            result += `/${id[0]}`
        }
    }catch (e) {
        console.error("ERROR FROM: [PathFromIdGenerator]")
        return null
    }

    return result
}