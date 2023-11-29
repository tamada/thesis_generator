const exts = {
    'format_latex': "tex",
    'format_html': "html",
    'format_markdown': "md",
    'format_word': "docx"
}

const clearForms = () => {
    ['familyname', 'thesis_id', 'thesis_title', 'student_id', 'author', 'email']
        .forEach(id => {
            document.getElementById(id).value = ""
    })
    document.getElementById('laboratory').selectedIndex = 0
    document.getElementById('warning').innerHTML = ""
}

const createParagraph = (text) => {
    const p = document.createElement('p')
    p.innerText = text
    return p
}

const putError = (iconId, message) => {
    document.getElementById(iconId).hidden = false
    document.getElementById('warning').appendChild(createParagraph(message))
}

const validateEmail = (email) => {
    const re = /\S+@\S+\.\S+/;
    if (email == "") {
        return
    }
    if(!re.test(email)) {
        putError('email_error', 'メールアドレスの形式が正しくありません．')
    } else {
        const userId = email.split('@')[0]
        if(/[gi]\d+/.test(userId)) {
            document.getElementById('student_id').value = userId.slice(1)
        }
    }
}

const validateFamilyName = (familyname) => {
    if (familyname.match(/^[^\x01-\x7E\xA1-\xDF]+$/)) {
        putError('familyname_error', '性は半角英数字で入力してください．')
    }
}

const hideAllErrors = () => {
    ['familyname_error', 'email_error'].forEach(id => {
        document.getElementById(id).hidden = true
    })
    document.getElementById('generate').disabled = true
    document.getElementById('warning').innerHTML = ""
}

const isAllEntriesFilled = () => {
    const list = ['familyname', 'thesis_id', 'thesis_title', 'student_id', 'author', 'email']
    for(let i = 0; i < list.length; i++) {
        if(document.getElementById(list[i]).value == "") {
            return false
        }
    }
    if(document.getElementById('laboratory').selectedIndex == 0) {
        return false
    }
    return document.getElementById('warning').innerHTML == ""
}

const validate = () => {
    hideAllErrors()
    validateEmail(document.getElementById('email').value)
    validateFamilyName(document.getElementById('familyname').value)
    if(isAllEntriesFilled()) {
        document.getElementById('generate').disabled = false
    }
    return document.getElementById('warning').innerHTML == ""
}

const updateThesisId = () => {
    if(!validate()) 
        return
    const list = document.getElementsByName('format');
    let ext
    for(let i = 0; i < list.length; i++) {
        if(list[i].checked) {
            ext = exts[list[i].id];
            break;
        }
    }
    const year = document.getElementById('thesis_year').value;
    const name = document.getElementById('familyname').value;
    const degree = "b"
    if(document.getElementById('master').checked) {
        degree = "m"
    }
    document.getElementById('thesis_id').value = year + degree + "thesis_" + name + "." + ext;
}

const generateThesis = () => {
    const thesisId = document.getElementById('thesis_id').value;
    const authorName = document.getElementById('author').value;
    const email = document.getElementById('email').value;
    const year = document.getElementById('thesis_year').value;
    const supervisor = document.getElementById('laboratory').selectedOptions[0].innerText
    const degree = document.getElementById('bachelor').checked ? "bachelor" : "master"
    const format = document.querySelector('input[name="format"]:checked').id
}