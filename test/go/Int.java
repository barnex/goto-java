package go;

public class Int extends IntPtr{

	public int value;

	public Int(){
	}

	public Int(int value){
		this.value = value;
	}

	public int value(){
		return value;
	}

	public void set(int value){
		this.value = value;
	}

	public IntPtr addr(){
		return this;
	}

	public boolean equals(Object other){
		return other instanceof Int 
			&& ((Int)other).value == this.value;	
	}
}
