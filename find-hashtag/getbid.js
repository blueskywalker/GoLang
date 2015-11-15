var data=db.brandref.find()
data.forEach(function (obj) { print ( obj.bid +"\t"+obj.Name ) })
